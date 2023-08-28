package scopes

import (
	"fmt"
	"strings"

	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	"gorm.io/gorm"
)

func PostIndexedOnly() QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(
			"id", "title",
			"created_at", "published_at",
			"updated_at", "user_id",
			"status",
		)
	}
}

func PostDeepMatch(keyword string, status values.PostStatus) QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		wrappedKeyword := fmt.Sprintf("%%%s%%", strings.ToLower(keyword))
		dbs := db.Table("posts").
			Select("distinct posts.id",
				"posts.title", "posts.created_at",
				"posts.status", "posts.updated_at",
				"posts.published_at").
			Joins("full join post_tags on posts.id = post_tags.post_id").
			Joins("left join tags on post_tags.tag_id = tags.id")

		dbs = dbs.
			Where(
				dbs.Where("lower(posts.title) like ?", wrappedKeyword).
					Or("lower(posts.content) like ?", wrappedKeyword).
					Or("lower(tags.title) like ?", wrappedKeyword),
			)

		if status != values.None {
			dbs = dbs.Where("posts.status = ?", status)
		}

		return dbs
	}
}
