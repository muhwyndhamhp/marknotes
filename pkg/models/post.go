package models

import (
	"context"
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title          string
	Abstract       string
	Content        string
	EncodedContent template.HTML
	Status         values.PostStatus `gorm:"index"`
	PublishedAt    time.Time
	Slug           string                 `gorm:"index"`
	FormMeta       map[string]interface{} `gorm:"-"`
	UserID         uint
	User           User
	Tags           []*Tag `gorm:"many2many:post_tags;"`
	Comments       []*Comment
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
	status values.PostStatus,
	sortQuery string,
	keyword string,
) {
	qs := url.Values{}
	qs.Add("page", strconv.Itoa(page))

	if sortQuery != "" {
		qs.Add("sortBy", sortQuery)
	}

	if status != values.None {
		qs.Add("status", string(values.Published))
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
