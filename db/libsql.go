package db

import (
	"fmt"

	"github.com/muhwyndhamhp/marknotes/config"
	libsql "github.com/renxzen/gorm-libsql"
	"gorm.io/gorm"
)

var sqldb *gorm.DB

func init() {
	url := config.Get("LIBSQL_URL")
	auth := config.Get("LIBSQL_AUTH_TOKEN")

	db, err := gorm.Open(libsql.Open(fmt.Sprintf("%s?authToken=%s", url, auth)), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqldb = db.Debug()
}

func GetLibSQLDB() *gorm.DB {
	return sqldb
}
