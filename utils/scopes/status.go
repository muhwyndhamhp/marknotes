package scopes

import (
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"gorm.io/gorm"
)

type QueryScope func(db *gorm.DB) *gorm.DB

func WithStatus(status models.PostStatus) QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		if status == models.None {
			return db
		}
		return db.Where("status = ?", status)
	}
}
