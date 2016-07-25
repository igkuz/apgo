package apgo

import (
    "github.com/jinzhu/gorm"
)

type AppContext struct {
  DB        *gorm.DB
  Config    *APConfig
}
