package models

import (
	"context"
	"html/template"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title          string
	Content        string
	EncodedContent template.HTML
	Status         PostStatus
	FormMeta       map[string]interface{} `gorm:"-"`
}

type PostStatus string

const (
	None      PostStatus = ""
	Draft     PostStatus = "draft"
	Published PostStatus = "published"
)

type PostRepository interface {
	Upsert(ctx context.Context, value *Post) error
	GetByID(ctx context.Context, id uint) (*Post, error)
	Get(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) ([]Post, error)
	Delete(ctx context.Context, id uint) error
}

func (m *Post) AppendFormMeta(page int, onlyPublished bool) {
	m.FormMeta = map[string]interface{}{
		"IsLastItem":    true,
		"Page":          page,
		"PublishedOnly": onlyPublished,
	}
}
