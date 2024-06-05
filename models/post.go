package models

import (
	"html/template"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title           string
	Abstract        string
	HeaderImageURL  string
	Content         string
	EncodedContent  template.HTML
	MarkdownContent string
	Status          PostStatus `gorm:"index"`
	PublishedAt     time.Time
	Slug            string                 `gorm:"index"`
	FormMeta        map[string]interface{} `gorm:"-"`
	UserID          uint
	User            User
	Tags            []Tag `gorm:"many2many:post_tags;"`
	TagsLiteral     string
	Comments        []Comment
}
type PostStatus string

const (
	None      PostStatus = ""
	Draft     PostStatus = "draft"
	Published PostStatus = "published"
)
