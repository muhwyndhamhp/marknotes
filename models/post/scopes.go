package post

import (
	"fmt"
	"strings"

	"github.com/muhwyndhamhp/marknotes/models"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"gorm.io/gorm"
)

func WithStatus(status models.PostStatus) scopes.QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("posts.status = ?", status)
	}
}

func Shallow() scopes.QueryScope {
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

func WithKeyword(keyword string) scopes.QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		wrappedKeyword := fmt.Sprintf("%%%s%%", strings.ToLower(keyword))
		dbs := db.Table("posts").
			Select(
				"distinct posts.id",
				"posts.title",
				"posts.created_at",
				"posts.status",
				"posts.updated_at",
				"posts.published_at",
			).
			Joins("full join post_tags on posts.id = post_tags.post_id").
			Joins("left join tags on post_tags.tag_id = tags.id")

		dbs = dbs.
			Where(
				dbs.Where("lower(posts.title) like ?", wrappedKeyword).
					Or("lower(posts.content) like ?", wrappedKeyword).
					Or("lower(tags.title) like ?", wrappedKeyword),
			)

		return dbs
	}
}
