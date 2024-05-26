package migration

import (
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Post{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.Tag{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.Comment{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Commenter{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Analytics{})
	if err != nil {
		panic(err)
	}
}
