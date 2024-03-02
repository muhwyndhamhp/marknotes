package db

import (
	"gorm.io/gorm"
)

var database *gorm.DB

func init() {
	// dsn := config.Get("DATABASE_URL")
	//
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	panic(err)
	// }
	// database = db
}

func GetDB() *gorm.DB {
	return database
}
