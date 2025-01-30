package internal

import (
	"context"
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"gorm.io/gorm"
)

type PostStatus string

const (
	PostStatusNone      PostStatus = ""
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
)

type Post struct {
	gorm.Model
	Title           string
	Abstract        string
	HeaderImageURL  string
	Content         string
	EncodedContent  template.HTML
	MarkdownContent string
	Status          PostStatus `gorm:"index"`
	PublishedAt     time.Time
	Slug            string                 `gorm:"index"`
	FormMeta        map[string]interface{} `gorm:"-"`
	UserID          uint
	User            User
	Tags            []*Tag `gorm:"many2many:post_tags;"`
	TagsLiteral     string
	Comments        []*Comment
}

type PostRepository interface {
	Upsert(ctx context.Context, value *Post) error
	GetByID(ctx context.Context, id uint) (*Post, error)
	Get(ctx context.Context, funcs ...scopes.QueryScope) ([]Post, error)
	Count(ctx context.Context, funs ...scopes.QueryScope) int
	Delete(ctx context.Context, id uint) error
}

func (m *Post) GenerateURL() string {
	baseURL := strings.Split(config.Get(config.OAUTH_URL), "/callback")[0]
	return fmt.Sprintf("%s/%s.html", baseURL, m.Slug)
}

func (m *Post) AppendFormMeta(
	page int,
	status PostStatus,
	sortQuery string,
	keyword string,
) {
	qs := url.Values{}
	qs.Add("page", strconv.Itoa(page))

	if sortQuery != "" {
		qs.Add("sortBy", sortQuery)
	}

	if status != PostStatusNone {
		qs.Add("status", string(PostStatusPublished))
	}

	if keyword != "" {
		qs.Add("search", keyword)
	}

	m.FormMeta = map[string]interface{}{
		"IsLastItem": true,
		"NextPath":   fmt.Sprintf("/posts?%s", qs.Encode()),
	}

	if keyword != "" {
		m.FormMeta["Keyword"] = keyword
	}
}

func WithStatus(status PostStatus) scopes.QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		if status == PostStatusNone {
			return db
		}
		return db.Where("status = ?", status)
	}
}

func PostDeepMatch(keyword string, status PostStatus) scopes.QueryScope {
	return func(db *gorm.DB) *gorm.DB {
		wrappedKeyword := fmt.Sprintf("%%%s%%", strings.ToLower(keyword))
		dbs := db.Table("posts").
			Select(
				"distinct posts.id",
				"posts.title",
				"posts.created_at",
				"posts.status",
				"posts.updated_at",
				"posts.published_at",
			).
			Joins("full join post_tags on posts.id = post_tags.post_id").
			Joins("left join tags on post_tags.tag_id = tags.id")

		dbs = dbs.
			Where(
				dbs.Where("lower(posts.title) like ?", wrappedKeyword).
					Or("lower(posts.content) like ?", wrappedKeyword).
					Or("lower(tags.title) like ?", wrappedKeyword),
			)

		if status != PostStatusNone {
			dbs = dbs.Where("posts.status = ?", status)
		}

		return dbs
	}
}
