package main

import (
	m "../models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

func main() {
	db, err := gorm.Open("mysql", "root:123@/sgap_development")
	if err != nil {
		log.Println("Error while connecting to DB: ", err)
	} else {
		db.AutoMigrate(m.Account{}, m.Ticket{})
		log.Println("Migration complete.")

		acc := m.Account{
			ExtId:          1,
			Name:           "Test account",
			GaAccessToken:  "123",
			GaViewId:       "1",
			GaQuotaUser:    "1",
			GaRefreshToken: "456",
		}
		db.Create(&acc)

		log.Println("Account has been created.")
	}
}
