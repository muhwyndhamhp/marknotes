package internal

import (
	"context"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"

	"gorm.io/gorm"
)

type Reply struct {
	gorm.Model
	Moderation

	ArticleID uint
	Message   string
	Alias     string

	ParentID *uint

	Parent  *Reply  `gorm:"foreignKey:ParentID"`
	Replies []Reply `gorm:"foreignKey:ParentID"`
	Article Post    `gorm:"foreignKey:ArticleID"`

	EnableReply bool `gorm:"-:all"`
	Highlight   bool `gorm:"-:all"`
	Page        int  `gorm:"-:all"`

	HidePublicity bool
}

func (r *Reply) AfterFind(tx *gorm.DB) (err error) {
	r.EnableReply = true

	return nil
}

type ReplyRepository interface {
	FetchArticleReplies(ctx context.Context, articleID uint) ([]Reply, error)
	Fetch(ctx context.Context, scopes ...scopes.QueryScope) ([]Reply, int, error)
	CreateReply(ctx context.Context, reply *Reply) error
	HideReply(ctx context.Context, id uint) error
	UpdateModeration(ctx context.Context, id uint, mod Moderation) error
}
