package db

import (
	"github.com/muhwyndhamhp/marknotes/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB

func init() {
	dsn := config.Get("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database = db
}

func GetDB() *gorm.DB {
	return database
}
