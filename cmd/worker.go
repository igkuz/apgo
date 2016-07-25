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

type appContext struct {
  DB        *gorm.DB
  config    *apgo.APConfig
}

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
  for i:=1; i<=3; i++ {
    context.Wg.Add(1)
    go apgo.Worker(context, i, updates, done)
  }
  time.Sleep(time.Second * 2)
  close(updates)
  close(done)
  context.Wg.Wait()
}
