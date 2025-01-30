package post

import (
	"context"
	"fmt"
	"github.com/muhwyndhamhp/marknotes/internal"

	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// Count implements models.PostRepository.
func (r *repository) Count(ctx context.Context, funcs ...scopes.QueryScope) int {
	s := scopes.Unwrap(funcs...)
	count := int64(0)
	if err := r.db.WithContext(ctx).
		Model(&internal.Post{}).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		Scopes(s...).
		Count(&count).
		Error; err != nil {
		fmt.Println(err)
		return 0
	}
	return int(count)
}

// Delete implements models.PostRepository.
func (r *repository) Delete(ctx context.Context, id uint) error {
	if err := r.db.Delete(&internal.Post{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) Get(ctx context.Context, funcs ...scopes.QueryScope) ([]internal.Post, error) {
	var res []internal.Post
	s := scopes.Unwrap(funcs...)
	if err := r.db.WithContext(ctx).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		Scopes(s...).
		Find(&res).
		Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*internal.Post, error) {
	var res internal.Post
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

func (r *repository) Upsert(ctx context.Context, value *internal.Post) error {
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

func NewPostRepository(db *gorm.DB) internal.PostRepository {
	return &repository{
		db: db,
	}
}
