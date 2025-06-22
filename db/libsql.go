package db

import (
	"context"
	"fmt"
	"github.com/glebarez/sqlite"
	_ "github.com/glebarez/sqlite"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/db/migration"
	"github.com/muhwyndhamhp/marknotes/utils/cloudbucket"
	"github.com/muhwyndhamhp/marknotes/utils/imageprocessing"
	"gorm.io/gorm"
	"time"
)

var sqldb *gorm.DB

func init() {
	bucket := cloudbucket.NewS3Client(imageprocessing.NewClient())

	dbName := "marknotes.db"

	if config.Get(config.ENV) == "dev" {
		dbName = "mk-test-1.db"
	}

	if err := bucket.DownloadDB(context.Background(), dbName); err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open(fmt.Sprint("file:./", dbName)), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqldb = db

	migration.Migrate(db)

	go func() {
		for {
			BackupDB()
			time.Sleep(time.Minute * 10)
		}
	}()
}

func BackupDB() {
	dbName := "marknotes.db"

	if config.Get(config.ENV) == "dev" {
		dbName = "mk-test-1.db"
	}

	bucket := cloudbucket.NewS3Client(imageprocessing.NewClient())
	ctx := context.Background()
	if err := bucket.UploadDB(ctx, dbName); err != nil {
		fmt.Println("*** Failed to backup DB: ", err.Error())
	}
}

func GetLibSQLDB() *gorm.DB {
	return sqldb
}
