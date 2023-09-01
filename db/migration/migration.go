package migration

import (
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.Debug().AutoMigrate(&models.Post{})
	db.Debug().AutoMigrate(&models.User{})
	db.Debug().AutoMigrate(&models.Tag{})
}
