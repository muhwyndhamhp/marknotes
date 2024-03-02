package models

import (
	"context"

	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string
	Name        string
	OauthUserID string
	Posts       []Post
}

type UserRepository interface {
	Upsert(ctx context.Context, value *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByOauthID(ctx context.Context, id string) (*User, error)
	Get(ctx context.Context, funcs ...scopes.QueryScope) ([]User, error)
	Delete(ctx context.Context, id uint) error
}
