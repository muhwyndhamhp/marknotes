package tag

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/pkg/tag/dto"
	pub_postlist "github.com/muhwyndhamhp/marknotes/pub/components/postlist"
	"github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
)

type TagFrontend struct {
	repo models.TagRepository
}

func NewTagFrontend(g *echo.Group, repo models.TagRepository, authMid echo.MiddlewareFunc) {
	fe := &TagFrontend{
		repo: repo,
	}

	g.POST("/tags/find-or-create", fe.TagsFindOrCreate, authMid)
	g.GET("/tags/remove", fe.TagsRemove, authMid)
}

func (ge *TagFrontend) TagsRemove(c echo.Context) error {
	return c.HTML(http.StatusOK, "")
}

func (fe *TagFrontend) TagsFindOrCreate(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.TagFindOrCreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	tagName := strings.ToLower(strings.TrimSpace(req.Tag))
	tagSlug := strings.ReplaceAll(tagName, " ", "-")

	tags, err := fe.repo.Get(ctx, scopes.Where("slug = ?", tagSlug))
	if err != nil {
		return errs.Wrap(err)
	}

	var tag models.Tag

	if len(tags) <= 0 {
		tag = models.Tag{
			Slug:  tagSlug,
			Title: strings.Title(tagName),
		}

		if err := fe.repo.Upsert(ctx, &tag); err != nil {
			return errs.Wrap(err)
		}
	} else {
		tag = tags[0]
	}

	models.SetTagEditable(&tag)

	for i := range req.Tags {
		id, _ := strconv.Atoi(req.Tags[i])
		if tag.ID == uint(id) {
			return c.HTML(http.StatusOK, "")
		}
	}
	tagItem := pub_postlist.TagItem(&tag)

	return template.AssertRender(c, http.StatusOK, tagItem)
}
