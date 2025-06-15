package dashboard

import (
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/internal"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/common/variables"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/dashboard/comment"
	templateRenderer "github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"github.com/muhwyndhamhp/marknotes/utils/tern"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
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

	rs, err := fe.getComments(ctx, 0, nil, "")
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

	rs, err := fe.getComments(ctx, 0, nil, "")
	if err != nil {
		return err
	}

	cm := comment.CommentsBody(rs)

	return templateRenderer.AssertRender(c, http.StatusOK, cm)
}

func (fe *handler) Comments(c echo.Context) error {
	ctx := c.Request().Context()

	f := map[string]string{}

	err := c.Bind(&f)
	if err != nil {
		return err
	}

	page, _ := strconv.Atoi(f["page"])

	var status *internal.ModerationStatus
	if v, ok := f["moderationStatus"]; ok {
		switch v {
		case internal.ModerationNothing.String():
			status = tern.ToPointer(internal.ModerationNothing)
		case internal.ModerationUnverified.String():
			status = tern.ToPointer(internal.ModerationUnverified)
		case internal.ModerationOK.String():
			status = tern.ToPointer(internal.ModerationOK)
		case internal.ModerationWarning.String():
			status = tern.ToPointer(internal.ModerationWarning)
		case internal.ModerationDangerous.String():
			status = tern.ToPointer(internal.ModerationDangerous)
		}
	}

	search := ""
	withSearch := false
	if v, ok := f["search"]; ok {
		search = v
		withSearch = ok
	}

	rs, err := fe.getComments(ctx, page, status, search)
	if err != nil {
		return err
	}

	var cm templ.Component
	if page == 0 && status == nil && !withSearch {
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

func (fe *handler) getComments(ctx context.Context, page int, status *internal.ModerationStatus, search string) ([]internal.Reply, error) {
	if page == 0 {
		page++
	}

	sc := []scopes.QueryScope{
		scopes.OrderBy("created_at", scopes.Descending),
		scopes.Where("hide_publicity != true"),
		scopes.Preload("Replies", func(db *gorm.DB) *gorm.DB { return db.Order("created_at DESC") }),
		scopes.Preload("Parent"),
		scopes.Preload("Article", func(db *gorm.DB) *gorm.DB { return db.Select("id", "title", "status") }),
		scopes.Paginate(page, 5),
	}

	if status != nil && *status != internal.ModerationNothing {
		sc = append(sc, scopes.Where("moderation_status = ?", *status))
	}

	if search != "" {
		sc = append(sc, scopes.Where(fmt.Sprint(
			"(LOWER(message) LIKE '%",
			strings.ToLower(search),
			"%' OR LOWER(alias) LIKE '%",
			strings.ToLower(search),
			"%')",
		)))
	}

	replies, _, err := fe.App.ReplyRepository.Fetch(
		ctx,
		sc...,
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
		item.Search = search

		if status != nil {
			item.FilterModerationStatus = status.String()
		}

		return item
	})

	return rs, nil
}
