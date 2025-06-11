package reply

import (
	"context"
	"errors"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"

	"github.com/muhwyndhamhp/marknotes/internal"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) UpdateModeration(ctx context.Context, id uint, mod internal.Moderation) error {
	return r.db.WithContext(ctx).
		Model(&internal.Reply{}).
		Where("id = ?", id).
		Updates(&internal.Reply{Moderation: mod}).Error
}

func (r *repository) HideReply(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&internal.Reply{}).
		Where("id = ?", id).
		Update("hide_publicity", true).Error
}

func (r *repository) Fetch(ctx context.Context, s ...scopes.QueryScope) ([]internal.Reply, int, error) {
	var res []internal.Reply
	var count int64

	q := r.db.WithContext(ctx).
		Model(&internal.Reply{}).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		Scopes(scopes.Unwrap(s...)...)

	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Find(&res).Error; err != nil {
		return nil, 0, err
	}

	return res, int(count), nil
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
