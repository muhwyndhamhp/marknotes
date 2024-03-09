package models

import (
	"context"

	"gorm.io/gorm"
)

type Commenter struct {
	gorm.Model
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Comments []*Comment `json:"comments"`
	// is blocked default false
	IsBlocked bool `json:"is_blocked" gorm:"default:false"`
}

type Comment struct {
	gorm.Model
	Text        string     `json:"text"`
	CommenterID uint       `json:"commenter_id"`
	Commenter   *Commenter `json:"commenter"`
	PostID      uint       `json:"post_id"`
}

type CommenterRepository interface {
	CreateCommenter(ctx context.Context, commenter *Commenter) error
	BlockCommenter(ctx context.Context, commenterID uint) error
	CreateComment(ctx context.Context, comment *Comment) error
	DeleteComment(ctx context.Context, commentID uint) error
	FindCommentByPostID(ctx context.Context, postID int) ([]Comment, error)
}
