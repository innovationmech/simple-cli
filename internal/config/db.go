package config

import (
	"sync"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open(sqlite.Open(viper.GetString("db.url")), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	})
	return db
}
