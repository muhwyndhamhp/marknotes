package scopes

import (
	"fmt"
	"time"

	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	"gorm.io/gorm"
)

func WithID(id uint) QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func Preload(query string, args ...interface{}) QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(query, args...)
	}
}

func Where(statement string, params ...interface{}) QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(statement, params...)
	}
}

func WithStatus(status values.PostStatus) QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		if status == values.None {
			return db
		}
		return db.Where("status = ?", status)
	}
}

func Between(field string, floor, ceil interface{}) QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), floor, ceil)
	}
}

func CreatedAfter(ct *time.Time) QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at > ?", ct)
	}
}

func CreatedBefore(ct *time.Time) QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at <= ?", ct)
	}
}

func UpdatedAfter(ct *time.Time) QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("updated_at > ?", ct)
	}
}

func UpdatedBefore(ct *time.Time) QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("updated_at <= ?", ct)
	}
}
