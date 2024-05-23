package main

import (
	"github.com/muhwyndhamhp/marknotes/db"
	"github.com/muhwyndhamhp/marknotes/db/migration"
)

func main() {
	migration.Migrate(db.GetLibSQLDB())
}
