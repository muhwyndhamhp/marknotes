package internal

import (
	"context"

	"gorm.io/gorm"
)

type Reply struct {
	gorm.Model
	Moderation

	ArticleID uint
	Message   string
	Alias     string

	ParentID *uint
	Replies  []Reply `gorm:"foreignKey:ParentID"`
}

type ReplyRepository interface {
	FetchArticleReplies(ctx context.Context, articleID uint) ([]Reply, error)
	CreateReply(ctx context.Context, reply *Reply) error
}
