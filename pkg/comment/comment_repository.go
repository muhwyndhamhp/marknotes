package comment

import (
	"context"

	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) BlockCommenter(ctx context.Context, commenterID uint) error {
	err := r.db.
		Model(&models.Commenter{}).
		Where("commenter_id =?", commenterID).
		Update("is_blocked", true).
		Error
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (r *repository) DeleteComment(ctx context.Context, commentID uint) error {
	err := r.db.Delete(&models.Comment{}, commentID).Error
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (r *repository) CreateComment(ctx context.Context, comment *models.Comment) error {
	err := r.db.Create(comment).Error
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (r *repository) CreateCommenter(ctx context.Context, commenter *models.Commenter) error {
	err := r.db.Create(commenter).Error
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (r *repository) FindCommentByPostID(ctx context.Context, postID int) ([]models.Comment, error) {
	var res []models.Comment
	err := r.db.Where("post_id =?", postID).Find(&res).Error
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return res, nil
}

func NewCommentRepository(db *gorm.DB) models.CommenterRepository {
	return &repository{
		db: db,
	}
}
