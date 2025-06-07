package migration

import (
	"fmt"
	"github.com/muhwyndhamhp/marknotes/internal"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	fmt.Println("** Start Migrate")
	err := db.AutoMigrate(&internal.User{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&internal.Post{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&internal.Tag{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&internal.Analytics{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&internal.Reply{})
	if err != nil {
		panic(err)
	}
}
