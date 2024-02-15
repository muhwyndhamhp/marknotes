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

// Delete implements models.TagRepository.
func (r *repository) Delete(ctx context.Context, id uint) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

// Get implements models.TagRepository.
func (r *repository) Get(ctx context.Context, funcs ...scopes.QueryScope) ([]models.Tag, error) {
	var res []models.Tag
	scopes := scopes.Unwrap(funcs...)
	err := r.db.WithContext(ctx).Debug().
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		Scopes(scopes...).
		Find(&res).
		Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetByID implements models.TagRepository.
func (r *repository) GetByID(ctx context.Context, id uint) (*models.Tag, error) {
	var res models.Tag
	if err := r.db.WithContext(ctx).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		First(&res, id).
		Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return &res, nil
}

// Upsert implements models.TagRepository.
func (r *repository) Upsert(ctx context.Context, value *models.Tag) error {
	if err := r.db.WithContext(ctx).Save(value).Error; err != nil {
		return err
	}
	return nil
}

func NewTagRepository(db *gorm.DB) models.TagRepository {
	return &repository{
		db: db,
	}
}
