package main

import (
  "github.com/igkuz/apgo"
  //m "github.com/igkuz/apgo/models"

  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"

  "log"
  "os"
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

  context := &appContext{
      DB: db,
      config: config,
  }
}
