package main

import (
	"github.com/muhwyndhamhp/marknotes/db"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
)

func main() {
	db := db.GetDB()

	db.Debug()

	db.AutoMigrate(&models.Post{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Tag{})
}
