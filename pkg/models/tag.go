package models

import (
	"context"

	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Slug     string
	Title    string
	Posts    []*Post                `gorm:"many2many:post_tags;"`
	FormMeta map[string]interface{} `gorm:"-"`
}

type TagRepository interface {
	Upsert(ctx context.Context, value *Tag) error
	GetByID(ctx context.Context, id uint) (*Tag, error)
	Get(ctx context.Context, funcs ...scopes.QueryScope) ([]Tag, error)
	Delete(ctx context.Context, id uint) error
}

func SetTagEditable(tags ...*Tag) {
	for i := range tags {
		tags[i].FormMeta = map[string]interface{}{
			"IsEditable": true,
		}
	}
}
