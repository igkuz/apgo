package main

import (
  "time"
  "sync"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

type appContext struct {
  db              *gorm.DB
  //globDayLimit    *Limiter
  //glob100SecLimit    *Limiter
}

type TicketUpdate struct {
  ID          int
  AccID       int
  PageViews   int
}

var wg sync.WaitGroup

func worker(ID int, updatesChannel chan<- *TicketUpdate, doneChannel <-chan int, context *appContext) {
  log.Println("Worker for AccountID: " ID, " stared...")
  fiveMinTick := time.NewTicker(time.Minute * 5)
  log.Println("5 min ticker stared: " time.Now())
  //oneHourTick := time.NewTicker(time.Hour * 1)
  //oneDayTick := time.NewTicker(time.Day * 1)
  for {
    select {
      case <- fiveMinTick.C:
        log.Println("Worker: ", ID, " start processing 5 min queue...")
        var r []AccJoinTicket
        context.db.Table("accounts AS a").Select("a.ext_id AS acc_id, a.ga_access_token, a.ga_view_id, a.ga_quota_user, a.ga_refresh_token, t.ext_id, t.url, t.published_at, t.url").Joins("LEFT JOIN tickets AS t ON a.ext_id = t.account_ext_id AND t.published_at >= ?", time.Now().Format("2006-01-01")).Scan(&r)
      case <- doneChannel:
        defer(wg.Done())
        return
    }
  }
}

func processRequest(accID int, updates chan<- *TicketUpdate, done <- chan int) {
  //TODO: move to datastore and config with loading from ENV
  db.Table
}

func processUpdates(updatesChannel <-chan *TicketUpdate, doneChannel <-chan int) {
  ticker := time.NewTicker(time.Millisecond * 500)
  tickets = make(map[int][]*TicketUpdate)
  for {
    select {
      case tu <- updatesChannel:
        tickets = append(tickets[tu.AccID], tu)
      case <- ticker.C:
        log.Println("Sending updates to api: ", tickets)
        for k := range tickets {
          delete(tickets, k)
        }
        log.Println("Cleared tickets map: ", tickets)
        // send values to API
      case <- doneChannel:
        defer(wg.Done())
        return
    }
  }
}

func main() {
  done := make(chan int)
  updates := make(chan *TicketUpdate, 200)
  db, err := gorm.Open("mysql", "root:123@/sgap_development")
  if err != nil {
    log.Println("Error while connecting to DB: ", err)
    return
  }
  context := &appContext{db: db}
  // on recieve command createAccount we should create worker for it
  // accLimit, glob2SecLimit, globDayLimit
  wg.Add(1)
  go worker(accId, updates, done, accLimit, context)
  wg.Add(1)
  go processUpdates(updates, done)
  wg.Wait()
}
