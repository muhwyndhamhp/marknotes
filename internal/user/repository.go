package user

import (
	"context"
	"github.com/muhwyndhamhp/marknotes/internal"

	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

var userCache = map[string]internal.User{}

func (r *repository) GetCache(ctx context.Context, email string) *internal.User {
	u, ok := userCache[email]

	if !ok {
		usrs, _ := r.Get(ctx, scopes.Where("email = ?", email))
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
	if err := r.db.Delete(&internal.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) Get(ctx context.Context, funcs ...scopes.QueryScope) ([]internal.User, error) {
	var res []internal.User
	scope := scopes.Unwrap(funcs...)
	err := r.db.WithContext(ctx).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		Scopes(scope...).
		Find(&res).
		Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*internal.User, error) {
	var res internal.User
	if err := r.db.WithContext(ctx).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		First(&res, id).
		Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return &res, nil
}

func (r *repository) GetByOauthID(ctx context.Context, id string) (*internal.User, error) {
	var res internal.User
	if err := r.db.WithContext(ctx).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		Where("oauth_user_id = ?", id).
		First(&res).
		Error; err != nil {
		return nil, errs.Wrap(err)
	}

	return &res, nil
}

func (r *repository) Upsert(ctx context.Context, value *internal.User) error {
	if err := r.db.WithContext(ctx).Save(value).Error; err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *gorm.DB) internal.UserRepository {
	return &repository{
		db: db,
	}
}
