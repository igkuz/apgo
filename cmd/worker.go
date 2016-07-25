package main

import (
  "github.com/igkuz/apgo"
  m "github.com/igkuz/apgo/models"

  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"

  "log"
  "os"
  "sync"
)

func main() {
  if os.Getenv("APP_ENV") == "" {
    os.Setenv("APP_ENV", "development")
  }
  config := apgo.NewConfig()
  log.Println(config)
  log.Println(config.GetDbString())
  db, err := gorm.Open(config.DB["dialect"], config.GetDbString())
  if err != nil {
    log.Fatalf("Error while connecting to DB: %v", err)
  }
  var wg sync.WaitGroup

  context := &apgo.AppContext{
      DB: db,
      Config: config,
      Wg: wg,
  }

  updates := make(chan *m.TicketUpdate, 200)
  done := make(chan int)
  context.Wg.Add(1)
  go apgo.Worker(context, 1, updates, done)
  go apgo.ProcessUpdates(context, updates, done)
  time.Sleep(time.Minute * 3)
  close(updates)
  close(done)
  context.Wg.Wait()
}
