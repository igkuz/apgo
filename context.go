package apgo

import (
	"github.com/jinzhu/gorm"
	"sync"
)

type AppContext struct {
	DB     *gorm.DB
	Config *APConfig
	Wg     sync.WaitGroup
}
