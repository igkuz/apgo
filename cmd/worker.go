package main

import (
  "github.com/igkuz/apgo"
  m "github.com/igkuz/apgo/models"

  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"

  "log"
  "os"
  "os/signal"
  "sync"
  "syscall"
) 
func main() {
  if os.Getenv("APP_ENV") == "" {
    os.Setenv("APP_ENV", "development")
  }

  context := prepareContext()

  signals := make(chan os.Signal, 1)
  updates := make(chan *m.TicketUpdate, 200)
  done := make(chan int)
  signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

  context.Wg.Add(1)
  go apgo.Worker(context, 1, updates, done)
  context.Wg.Add(1)
  go apgo.ProcessUpdates(context, updates, done)

  for {
    select {
    case sig := <-signals:
      log.Printf("%v signal was sent, exiting", sig)
      close(updates)
      close(done)
      context.Wg.Wait()
    default:
    }
  }
  //time.Sleep(time.Minute * 3)
  //close(updates)
  //close(done)
  //context.Wg.Wait()
}

func prepareContext() *apgo.AppContext {
  config := apgo.NewConfig()
  db, err := gorm.Open(config.DB["dialect"], config.GetDbString())
  if err != nil {
    log.Fatalf("Error while connecting to DB: %v", err)
  }
  var wg sync.WaitGroup

  return &apgo.AppContext{
      DB: db,
      Config: config,
      Wg: wg,
  }
}
