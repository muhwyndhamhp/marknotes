package repository

import (
	"context"

	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// GetByOauthID implements models.UserRepository.

func (r *repository) Delete(ctx context.Context, id uint) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) Get(ctx context.Context, funcs ...scopes.QueryScope) ([]models.User, error) {
	var res []models.User
	scopes := scopes.Unwrap(funcs...)
	err := r.db.WithContext(ctx).
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
	if err := r.db.WithContext(ctx).First(&res, id).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return &res, nil
}

func (r *repository) GetByOauthID(ctx context.Context, id string) (*models.User, error) {
	var res models.User
	if err := r.db.WithContext(ctx).Where("oauth_user_id = ?", id).First(&res).Error; err != nil {
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
