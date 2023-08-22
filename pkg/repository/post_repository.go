package repository

import (
	"context"

	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// Delete implements models.PostRepository.
func (r *repository) Delete(ctx context.Context, id uint) error {
	if err := r.db.Delete(&models.Post{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) Get(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) ([]models.Post, error) {
	var res []models.Post
	err := r.db.WithContext(ctx).
		Scopes(funcs...).
		Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*models.Post, error) {
	var res models.Post
	if err := r.db.WithContext(ctx).First(&res, id).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return &res, nil
}

func (r *repository) Upsert(ctx context.Context, value *models.Post) error {
	if err := r.db.WithContext(ctx).Save(value).Error; err != nil {
		return err
	}
	return nil
}

func NewPostRepository(db *gorm.DB) models.PostRepository {
	return &repository{
		db: db,
	}
}
