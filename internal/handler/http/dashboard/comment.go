package dashboard

import (
	"context"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/internal"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/common/variables"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/dashboard/comment"
	templateRenderer "github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func (fe *handler) MarkSafe(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))

	now := time.Now()
	if err := fe.App.ReplyRepository.UpdateModeration(ctx, uint(id), internal.Moderation{
		LastModeratedAt:  &now,
		ModerationStatus: internal.ModerationOK,
	}); err != nil {
		return err
	}

	rs, err := fe.getComments(ctx, 0)
	if err != nil {
		return err
	}

	cm := comment.CommentsBody(rs)

	return templateRenderer.AssertRender(c, http.StatusOK, cm)
}

func (fe *handler) HideComment(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))

	if err := fe.App.ReplyRepository.HideReply(ctx, uint(id)); err != nil {
		return err
	}

	rs, err := fe.getComments(ctx, 0)
	if err != nil {
		return err
	}

	cm := comment.CommentsBody(rs)

	return templateRenderer.AssertRender(c, http.StatusOK, cm)
}

func (fe *handler) Comments(c echo.Context) error {
	ctx := c.Request().Context()

	page, _ := strconv.Atoi(c.QueryParam("page"))

	rs, err := fe.getComments(ctx, page)
	if err != nil {
		return err
	}

	var cm templ.Component
	if page == 0 {
		cm = comment.Comments(comment.CommentsVM{
			Opts: variables.DashboardOpts{
				Nav:         nav(1),
				BreadCrumbs: fe.Breadcrumbs("dashboard/comments"),
			},
			Comments: rs,
		})
	} else {
		cm = comment.CommentBody(rs)
	}

	return templateRenderer.AssertRender(c, http.StatusOK, cm)
}

func (fe *handler) getComments(ctx context.Context, page int) ([]internal.Reply, error) {
	if page == 0 {
		page++
	}

	replies, _, err := fe.App.ReplyRepository.Fetch(
		ctx,
		scopes.OrderBy("created_at", scopes.Descending),
		scopes.Where("hide_publicity != true"),
		scopes.Preload("Replies", func(db *gorm.DB) *gorm.DB { return db.Order("created_at DESC") }),
		scopes.Preload("Parent"),
		scopes.Preload("Article", func(db *gorm.DB) *gorm.DB { return db.Select("id", "title", "status") }),
		scopes.Paginate(page, 5),
	)
	if err != nil {
		return nil, err
	}

	rs := lo.Map(replies, func(item internal.Reply, index int) internal.Reply {
		item.EnableReply = false
		item.Highlight = true
		item.Replies = lo.Map(item.Replies, func(i internal.Reply, index int) internal.Reply {
			i.EnableReply = false
			return i
		})

		if item.Parent != nil {
			item.Parent.EnableReply = false
		}

		item.Page = page + 1

		return item
	})

	return rs, nil
}
