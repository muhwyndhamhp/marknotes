package commenter

import (
	"context"
	"github.com/muhwyndhamhp/marknotes/internal"

	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) BlockCommenter(ctx context.Context, commenterID uint) error {
	err := r.db.
		Model(&internal.Commenter{}).
		Where("commenter_id =?", commenterID).
		Update("is_blocked", true).
		Error
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (r *repository) DeleteComment(ctx context.Context, commentID uint) error {
	err := r.db.Delete(&internal.Comment{}, commentID).Error
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (r *repository) CreateComment(ctx context.Context, comment *internal.Comment) error {
	err := r.db.Create(comment).Error
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (r *repository) CreateCommenter(ctx context.Context, commenter *internal.Commenter) error {
	err := r.db.Create(commenter).Error
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (r *repository) FindCommentByPostID(ctx context.Context, postID int) ([]internal.Comment, error) {
	var res []internal.Comment
	err := r.db.Where("post_id =?", postID).Find(&res).Error
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return res, nil
}

func NewCommentRepository(db *gorm.DB) internal.CommenterRepository {
	return &repository{
		db: db,
	}
}
