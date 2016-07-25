package apgo

import (
	m "github.com/igkuz/apgo/models"

	"github.com/jinzhu/gorm"

	"log"
	"time"
  "golang.org/x/time/rate"
)

func getTicketsForToday(db *gorm.DB) *gorm.DB {
	return db.Table("accounts AS a").
		Select(getSelectString()).
		Joins("LEFT JOIN tickets AS t ON a.ext_id = t.account_ext_id AND t.published_at >= ?", time.Now().Format("2006-01-01"))
}

func getSelectString() string {
	return "a.ext_id AS acc_id, a.ga_access_token, a.ga_view_id, a.ga_quota_user, a.ga_refresh_token, t.ext_id, t.url, t.published_at, t.url, t.ga_page_views"
}

func Worker(context *AppContext, ID int, updates chan<- *m.TicketUpdate, done <-chan int) {
  defer context.Wg.Done()
	lf := log.Printf

	lf("Worker for AccountID: %v started.\n", ID)
	// just for tests leave 1 minute
	fiveMinTick := time.NewTicker(time.Minute * 1)
	lf("5 min ticker started at: %v.\n", time.Now())
  accountLimit := getLimiter()

	for {
		select {
		case <-fiveMinTick.C:
			lf("Worker: %v started processing 5 min queue (today tickets).", ID)
			var result []m.AccJoinTicket
			context.DB.Scopes(getTicketsForToday).Scan(&result)
			log.Printf("Selected tickets: %#v", result)
			for _, t := range result {
				log.Println("Processing ticket: ", t)
        context.Wg.Add(1)
        go processAnalytics(context, &t, updates, done, accountLimit)
			}
		case <-done:
			lf("Done signal catched. Stopping worker: %v.\n", ID)
      return
		}
	}
}

// Processing updates from API calls. When Google Analytics or other system provides response. We should send update to another system API.
// The updates channel has buffer of 200 records, so we can grab it out in some small amount of time decreasing requests to external API by 5 rps.
func ProcessUpdates(context *AppContext, updates <-chan *m.TicketUpdate, done <-chan int) {
  defer context.Wg.Done()

  ticker := time.NewTicker(time.Millisecond * 200)
  tickets := make(map[int][]*m.TicketUpdate)

  for {
    select {
    case tu := <-updates:
        tickets[tu.AccId] = append(tickets[tu.AccId], tu)
    case <- ticker.C:
        log.Println("Sending updates to api: ", tickets)
        for k := range tickets {
          delete(tickets, k)
        }
        log.Println("Cleared tickets map: ", tickets)
        // send values to API
    case <- done:
        return
    }
  }
}

func getLimiter() *rate.Limiter {
  return rate.NewLimiter(10, 1)
}

func processAnalytics(context *AppContext, ticket *m.AccJoinTicket, updates chan<- *m.TicketUpdate, done <-chan int, accLimit *rate.Limiter) {
  defer context.Wg.Done()
  ok := accLimit.Allow()

  tu := &m.TicketUpdate{
    Id: 0,
    AccId: 0,
    PageViews: 0,
  }

  if ok {
      log.Printf("Processing ticket %v.\n", ticket.ExtId)
      time.Sleep(time.Millisecond * 50)
      log.Printf("Sending updates.")
      tu.Id = ticket.ExtId
      tu.AccId = ticket.AccId
      tu.PageViews = ticket.GaPageViews
  } else {
    log.Printf("No slot for request. Scheduling event.")
    context.Wg.Add(1)
    go processAnalytics(context, ticket, updates, done, accLimit)
  }

  for {
    select {
    case updates<- tu:
      log.Printf("Updates for ticket %v, were successfully sent.\n", tu.Id)
      return
    case <-done:
      log.Printf("Done signal catched. Stopping processAnalytics for ticket: %v.\n", ticket.ExtId)
      return
    }
  }
}
