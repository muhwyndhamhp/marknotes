package db

import (
	"fmt"

	"github.com/muhwyndhamhp/gotes-mx/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB

func init() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		config.Get(config.DB_HOST),
		config.Get(config.DB_PORT),
		config.Get(config.DB_USER),
		config.Get(config.DB_NAME),
		config.Get(config.DB_PASSWORD),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database = db
}

func GetDB() *gorm.DB {
	return database
}
