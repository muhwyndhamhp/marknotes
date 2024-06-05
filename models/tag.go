package models

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Slug     string
	Title    string
	Posts    []Post                 `gorm:"many2many:post_tags;"`
	FormMeta map[string]interface{} `gorm:"-"`
}
