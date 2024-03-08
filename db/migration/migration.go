package migration

import (
	"github.com/muhwyndhamhp/marknotes/db"
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

func MigrateToLibSQL() {
	var usr models.User
	err := db.GetDB().First(&usr).Error
	if err != nil {
		panic(err)
	}

	err = db.GetLibSQLDB().Create(&usr).Error
	if err != nil {
		panic(err)
	}

	var posts []models.Post
	err = db.GetDB().Find(&posts).Error
	if err != nil {
		panic(err)
	}

	for _, post := range posts {
		err = db.GetLibSQLDB().Create(&post).Error
		if err != nil {
			panic(err)
		}
	}

	var tags []models.Tag
	err = db.GetDB().Find(&tags).Error
	if err != nil {
		panic(err)
	}

	for _, tag := range tags {
		err = db.GetLibSQLDB().Create(&tag).Error
		if err != nil {
			panic(err)
		}
	}
}
