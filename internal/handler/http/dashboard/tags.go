package dashboard

import (
	"encoding/json"
	"fmt"
	"github.com/muhwyndhamhp/marknotes/internal"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/dashboard/tags"
	templates "github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
)

func (fe *handler) Tags(c echo.Context) error {
	ctx := c.Request().Context()

	tagQuery := c.QueryParam("tag")

	tagName := strings.ToLower(strings.TrimSpace(tagQuery))
	tagSlug := strings.ReplaceAll(tagName, " ", "-")

	tl, err := fe.App.TagRepository.Get(
		ctx,
		scopes.Where("slug LIKE ?", fmt.Sprintf("%%%s%%", tagSlug)),
		scopes.Paginate(1, 5),
	)
	if err != nil {
		return err
	}
	if len(tl) == 0 {
		tl = append(tl, internal.Tag{
			Slug:  tagSlug,
			Title: tagQuery,
		})
	}

	var tagTs []string
	for i := range tl {
		tagTs = append(tagTs, tl[i].Title)
	}

	js, _ := json.Marshal(tagTs)
	c.Response().Header().Set("X-Tags", string(js))

	tagSuggest := tags.TagSuggest(tl)

	return templates.AssertRender(c, http.StatusOK, tagSuggest)
}
