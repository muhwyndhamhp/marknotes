package scopes

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

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
