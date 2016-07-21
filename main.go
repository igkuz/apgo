package main

import (
  "time"
  "sync"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)


type TicketUpdate struct {
  ID          int
  AccID       int
  PageViews   int
}

var wg sync.WaitGroup

func worker(ID int, updatesChannel chan<- *TicketUpdate, doneChannel <-chan int, limiter *Limiter) {
  log.Println("Worker for AccountID: " ID, " stared...")
  fiveMinTick := time.NewTicker(time.Minute * 5)
  log.Println("5 min ticker stared: " time.Now())
  //oneHourTick := time.NewTicker(time.Hour * 1)
  //oneDayTick := time.NewTicker(time.Day * 1)
  for {
    select {
      case <- fiveMinTick.C:
        log.Println("Worker: ", ID, " start processing 5 min queue...")
        // Select from DB new tickets
        go processRequest(ticket
      case <- doneChannel:
        defer(wg.Done())
        return
    }
  }
}

func processRequest(accID int, updates chan<- *TicketUpdate, done <- chan int) {
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
  // on recieve command createAccount we should create worker for it
  // accLimit, glob2SecLimit, globDayLimit
  wg.Add(1)
  go worker(accId, updates, done, accLimit, glob2SecLimit)
  wg.Add(1)
  go processUpdates(updates, done)
  wg.Wait()
}
