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

func (r *repository) Get(ctx context.Context, queryOpts scopes.QueryOpts) ([]models.Post, error) {
	var res []models.Post
	err := r.db.WithContext(ctx).
		Scopes(
			scopes.Paginate(queryOpts.Page, queryOpts.PageSize),
			scopes.OrderBy(queryOpts.Order, queryOpts.OrderDir),
		).
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
