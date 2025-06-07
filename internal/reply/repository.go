package reply

import (
	"context"
	"errors"

	"github.com/muhwyndhamhp/marknotes/internal"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) CreateReply(ctx context.Context, reply *internal.Reply) error {
	if reply.ArticleID == 0 {
		return errors.New("reply must have article id attached")
	}

	if reply.Alias == "" {
		return errors.New("reply must have user alias")
	}

	return r.db.
		WithContext(ctx).
		Save(reply).Error
}

func (r *repository) FetchArticleReplies(ctx context.Context, articleID uint) ([]internal.Reply, error) {
	var res []internal.Reply

	if err := r.db.WithContext(ctx).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		Where("moderation_status != ?", internal.ModerationDangerous).
		Where("article_id = ?", articleID).
		Preload("Replies", func(db *gorm.DB) *gorm.DB { return db.Order("created_at DESC") }).
		Order("created_at DESC").
		Find(&res).
		Error; err != nil {
		return nil, err
	}

	return res, nil
}

func NewRepository(db *gorm.DB) internal.ReplyRepository {
	return &repository{db: db}
}
