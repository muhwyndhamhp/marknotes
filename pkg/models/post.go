package models

import (
	"context"
	"html/template"

	"github.com/muhwyndhamhp/gotes-mx/utils/scopes"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title          string
	Content        string
	EncodedContent template.HTML
	Status         PostStatus
	FormMeta       map[string]FormMeta `gorm:"-"`
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
