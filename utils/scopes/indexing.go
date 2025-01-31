package scopes

import (
	"gorm.io/gorm"
)

func PostIndexedOnly() QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(
			"id", "title",
			"created_at", "published_at",
			"updated_at", "user_id",
			"status", "slug",
			"tags_literal",
		)
	}
}
