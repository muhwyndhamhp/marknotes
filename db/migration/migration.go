package main

import (
	"github.com/muhwyndhamhp/gotes-mx/db"
	"github.com/muhwyndhamhp/gotes-mx/pkg/models"
)

func main() {
	db := db.GetDB()

	db.Debug()

	db.AutoMigrate(&models.Post{})
}
