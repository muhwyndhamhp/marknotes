package scopes

import "gorm.io/gorm"

func ActiveJob() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status active")
	}
}
