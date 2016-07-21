package main

import (
	m "../models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
  "time"
)

func main() {
	db, err := gorm.Open("mysql", "root:123@/sgap_development")
	if err != nil {
		log.Println("Error while connecting to DB: ", err)
	} else {
		db.AutoMigrate(m.Account{}, m.Ticket{})
		log.Println("Migration complete.")

    tickets := []m.Ticket{}
    for i:=1; i <= 10; i++ {
      t := m.Ticket{
        ExtId: i,
        PublishedAt:   time.Now(),
        Url: "http://google.com",
        Active: true,
        GaPageViews: 0,
      }
      tickets = append(tickets, t)
    }
		acc := m.Account{
			ExtId:          1,
			Name:           "Test account",
			GaAccessToken:  "123",
			GaViewId:       "1",
			GaQuotaUser:    "1",
			GaRefreshToken: "456",
      Tickets: tickets,
		}
		db.Create(&acc)

		log.Println("Account has been created.")
	}
}
