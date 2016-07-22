package main

import (
  "time"
  "sync"
  "log"
  m "./models"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

type appContext struct {
  db              *gorm.DB
  //globDayLimit    *Limiter
  //glob100SecLimit    *Limiter
}

type TicketUpdate struct {
  Id          int
  AccId       int
  PageViews   int
}

var wg sync.WaitGroup

func worker(ID int, context *appContext, updatesChannel chan<- *TicketUpdate, doneChannel <-chan int) {
  ln := log.Println
  ln("Worker for AccountID: ", ID, " stared...")
  fiveMinTick := time.NewTicker(time.Second * 3)
  ln("5 min ticker stared: ", time.Now())
  //oneHourTick := time.NewTicker(time.Hour * 1)
  //oneDayTick := time.NewTicker(time.Day * 1)
  for {
    select {
      case <- fiveMinTick.C:
        log.Println("Worker: ", ID, " start processing 5 min queue...")
        var r []m.AccJoinTicket
        context.db.Table("accounts AS a").Select("a.ext_id AS acc_id, a.ga_access_token, a.ga_view_id, a.ga_quota_user, a.ga_refresh_token, t.ext_id, t.url, t.published_at, t.url, t.ga_page_views").Joins("LEFT JOIN tickets AS t ON a.ext_id = t.account_ext_id AND t.published_at >= ?", time.Now().Format("2006-01-01")).Scan(&r)
        log.Printf("Rows: %#v", r)
        for _, t := range r {
          log.Println("Processing ticket: ", t)
          wg.Add(1)
          go getAnalytics(&t, updatesChannel)
        }
      case <- doneChannel:
        defer wg.Done()
        return
    }
  }
}

func getAnalytics(ticket *m.AccJoinTicket, updatesChannel chan<- *TicketUpdate) {
    log.Println("Doing some job...")
    time.Sleep(time.Millisecond*100)
    tu := &TicketUpdate{
      Id: ticket.ExtId,
      AccId: ticket.AccId,
      PageViews: ticket.GaPageViews + 1,
    }
    log.Println("Ticket updated, old page views: ", ticket.GaPageViews, " new page views: ", tu.PageViews, " ticket update: ", tu)
    updatesChannel <- tu
    defer wg.Done()
}

func processUpdates(updatesChannel <-chan *TicketUpdate, doneChannel <-chan int) {
  ticker := time.NewTicker(time.Millisecond * 1000)
  tickets := make(map[int][]*TicketUpdate)
  for {
    select {
      case tu := <- updatesChannel:
        tickets[tu.AccId] = append(tickets[tu.AccId], tu)
      case <- ticker.C:
        log.Println("Sending updates to api: ", tickets)
        for k := range tickets {
          delete(tickets, k)
        }
        log.Println("Cleared tickets map: ", tickets)
        // send values to API
      case <- doneChannel:
        defer wg.Done()
        return
    }
  }
}

func main() {
  done := make(chan int)
  updates := make(chan *TicketUpdate, 200)
  db, err := gorm.Open("mysql", "root:123@/sgap_development?parseTime=true")
  if err != nil {
    log.Println("Error while connecting to DB: ", err)
    return
  }
  context := &appContext{db: db}
  // on recieve command createAccount we should create worker for it
  // accLimit, glob2SecLimit, globDayLimit
  wg.Add(1)
  go worker(1, context, updates, done)
  wg.Add(1)
  go processUpdates(updates, done)
  wg.Wait()
}
