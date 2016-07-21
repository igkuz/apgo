package models

import (
	"github.com/jinzhu/gorm"
	"time"
  "fmt"
)

type Account struct {
	gorm.Model

	ExtId          int
	Name           string
	GaAccessToken  string
	GaViewId       string
	GaQuotaUser    string
	GaRefreshToken string
	Active         bool     `sql:"DEFAULT:true"`
	Tickets         []Ticket `gorm:"ForeignKey:AccountExtId"`
}

type Ticket struct {
	gorm.Model

	AccountExtId int `gorm:"not null"`
	ExtId        int `gorm:"not null"`
	PublishedAt  time.Time
	Url          string `gorm:"type:text;not null"`
	Active       bool   `sql:"DEFAULT:true"`
	GaPageViews  int
}

type AccJoinTicket struct {
  AccId               int
  GaAccessToken       string
  GaViewId            string
  GaQuotaUser         string
  GaRefreshToken      string
  ExtId               int
  Url                 string
  PublishedAt         time.Time
}

func (ajt AccJoinTicket) String() string {
  
    return "AccJoinTicket: {" + 
        fmt.Sprintf("TicketId: %v, ", ajt.ExtId) +
        fmt.Sprintf("Url: %v, ", ajt.Url) +
        fmt.Sprintf("PublishedAt: %v, ", ajt.PublishedAt) +
        fmt.Sprintf("AccId: %v, ", ajt.AccId) +
        fmt.Sprintf("GaAccessToken: %v, ", ajt.GaAccessToken) +
        fmt.Sprintf("GaRefreshToken: %v, ", ajt.GaRefreshToken) +
        fmt.Sprintf("GaViewId: %v, ", ajt.GaViewId) +
        fmt.Sprintf("GaQuotaUser: %v ", ajt.GaQuotaUser) +
    "}\n"
}
