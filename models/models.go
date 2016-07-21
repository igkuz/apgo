package models

import (
	"github.com/jinzhu/gorm"
	"time"
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
	Ticket         []Ticket `gorm:"ForeignKey:AccountExtId"`
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
