package repository

import (
	"context"
	"github.com/muhwyndhamhp/marknotes/db"

	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

var userCache = map[string]models.User{}

func (r *repository) GetCache(ctx context.Context, email string) *models.User {
	u, ok := userCache[email]

	if !ok {
		usrs, _ := r.Get(ctx, db.Where("email = ?", email))
		for _, usr := range usrs {
			userCache[usr.Email] = usr
		}
		u, ok = userCache[email]
		if !ok {
			return nil
		}
	}

	return &u
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) Get(ctx context.Context, funcs ...db.QueryScope) ([]models.User, error) {
	var res []models.User
	scopes := db.Unwrap(funcs...)
	err := r.db.WithContext(ctx).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		Scopes(scopes...).
		Find(&res).
		Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var res models.User
	if err := r.db.WithContext(ctx).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		First(&res, id).
		Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return &res, nil
}

func (r *repository) GetByOauthID(ctx context.Context, id string) (*models.User, error) {
	var res models.User
	if err := r.db.WithContext(ctx).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		Where("oauth_user_id = ?", id).
		First(&res).
		Error; err != nil {
		return nil, errs.Wrap(err)
	}

	return &res, nil
}

func (r *repository) Upsert(ctx context.Context, value *models.User) error {
	if err := r.db.WithContext(ctx).Save(value).Error; err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *gorm.DB) models.UserRepository {
	return &repository{
		db: db,
	}
}
