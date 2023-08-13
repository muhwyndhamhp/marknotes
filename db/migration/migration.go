package main

import (
	"github.com/muhwyndhamhp/gotes-mx/db"
)

func main() {
	db := db.GetDB()

	db.Debug()
}
