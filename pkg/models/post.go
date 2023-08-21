package models

import (
	"context"
	"html/template"

	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title          string
	Content        string
	EncodedContent template.HTML
	Preview        template.HTML
	Status         PostStatus
	FormMeta       map[string]interface{} `gorm:"-"`
}

type PostStatus string

const (
	Draft     PostStatus = "draft"
	Published PostStatus = "published"
)

type PostRepository interface {
	Upsert(ctx context.Context, value *Post) error
	GetByID(ctx context.Context, id uint) (*Post, error)
	Get(ctx context.Context, queryOpts scopes.QueryOpts) ([]Post, error)
}

func (m *Post) AppendFormMeta(page int) {
	m.FormMeta = map[string]interface{}{
		"IsLastItem": true,
		"Page":       page,
	}
}
