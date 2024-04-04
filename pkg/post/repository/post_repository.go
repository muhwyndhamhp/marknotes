package repository

import (
	"context"
	"fmt"

	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// Count implements models.PostRepository.
func (r *repository) Count(ctx context.Context, funcs ...scopes.QueryScope) int {
	scopes := scopes.Unwrap(funcs...)
	count := int64(0)
	if err := r.db.WithContext(ctx).
		Model(&models.Post{}).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		Scopes(scopes...).
		Count(&count).
		Error; err != nil {
		fmt.Println(err)
		return 0
	}
	return int(count)
}

// Delete implements models.PostRepository.
func (r *repository) Delete(ctx context.Context, id uint) error {
	if err := r.db.Delete(&models.Post{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) Get(ctx context.Context, funcs ...scopes.QueryScope) ([]models.Post, error) {
	var res []models.Post
	scopes := scopes.Unwrap(funcs...)
	if err := r.db.WithContext(ctx).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		Scopes(scopes...).
		Find(&res).
		Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*models.Post, error) {
	var res models.Post
	if err := r.db.WithContext(ctx).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		Preload("Tags").
		Where("id = ?", id).
		First(&res).
		Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return &res, nil
}

func (r *repository) Upsert(ctx context.Context, value *models.Post) error {
	if trxErr := r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.db.
			WithContext(ctx).
			Save(value).Error; err != nil {
			return err
		}

		if len(value.Tags) <= 0 {
			return nil
		}

		if err := r.db.
			WithContext(ctx).
			Model(value).
			Association("Tags").
			Replace(value.Tags); err != nil {
			return err
		}
		return nil
	}); trxErr != nil {
		return trxErr
	}

	return nil
}

func NewPostRepository(db *gorm.DB) models.PostRepository {
	return &repository{
		db: db,
	}
}
