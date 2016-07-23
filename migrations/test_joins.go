package main

import (
  m "../models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
  "time"
)


func main() {
	db, err := gorm.Open("mysql", "root:123@/sgap_development?parseTime=true")
	if err != nil {
		log.Println("Error while connecting to DB: ", err)
    return
	}
  var r []m.AccJoinTicket
  db.Table("accounts AS a").Select("a.ext_id AS acc_id, a.ga_access_token, a.ga_view_id, a.ga_quota_user, a.ga_refresh_token, t.ext_id, t.url, t.published_at, t.url, t.ga_page_views").Joins("LEFT JOIN tickets AS t ON a.ext_id = t.account_ext_id AND t.published_at >= ?", time.Now().Format("2006-01-01")).Scan(&r)
  log.Printf("Joined results: %v", r)
}
