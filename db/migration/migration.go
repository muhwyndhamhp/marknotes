package migration

import (
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.Debug().AutoMigrate(&models.Post{})
	if err != nil {
		panic(err)
	}
	err = db.Debug().AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
	err = db.Debug().AutoMigrate(&models.Tag{})
	if err != nil {
		panic(err)
	}
}
