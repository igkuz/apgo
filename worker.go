package apgo

import (
	m "github.com/igkuz/apgo/models"

	"github.com/jinzhu/gorm"

	"log"
	"time"
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
        go processAnalytics(context, t, updates, done)
			}
		case <-done:
			lf("Done signal catched. Stopping worker: %v.\n", ID)
      return
		}
	}
}

func processAnalytics(context *AppContext, ticket *AccJoinTicket, updates chan<- *m.TicketUpdate, done <-chan int) {
  defer context.Wg.Done()

  for {
    select {
    case <-done:
      lf("Done signal catched. Stopping processAnalytics for ticket: %v.\n", ticket.ExtId)
      return
    }
  }
}
