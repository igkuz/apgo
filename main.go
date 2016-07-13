package main

import (
  "time"
  "sync"
)

type TicketUpdate struct {
  ID          int
  AccID       int
  PageViews   int
}

var wg sync.WaitGroup

func worker(updatesChannel chan<- *TicketUpdate, doneChannel <-chan int, limiter *Limiter) {
  for {
    select {
      case <- doneChannel:
        defer(wg.Done())
        return
    }
  }
}

func processUpdates(updatesChannel <-chan *TicketUpdate, doneChannel <-chan int) {
  ticker := time.NewTicker(time.Millisecond * 500)
  for {
    tickets = make(map[int][]*TicketUpdate)
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
  go worker(updates, done, limiter)
  wg.Add(1)
  go processUpdates(updates, done)
  wg.Wait()
}
