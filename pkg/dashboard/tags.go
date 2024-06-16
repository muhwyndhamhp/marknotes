package dashboard

import (
	"encoding/json"
	"fmt"
	"github.com/muhwyndhamhp/marknotes/db"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	pub_tagsuggest "github.com/muhwyndhamhp/marknotes/pub/components/tagsuggest"
	templates "github.com/muhwyndhamhp/marknotes/template"
)

func (fe *DashboardFrontend) Tags(c echo.Context) error {
	ctx := c.Request().Context()

	tagQuery := c.QueryParam("tag")

	tagName := strings.ToLower(strings.TrimSpace(tagQuery))
	tagSlug := strings.ReplaceAll(tagName, " ", "-")

	tags, err := fe.TagRepo.Get(
		ctx,
		db.Where("slug LIKE ?", fmt.Sprintf("%%%s%%", tagSlug)),
		db.Paginate(1, 5),
	)
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		tags = append(tags, models.Tag{
			Slug:  tagSlug,
			Title: tagQuery,
		})
	}

	var tagTs []string
	for i := range tags {
		tagTs = append(tagTs, tags[i].Title)
	}

	js, _ := json.Marshal(tagTs)
	c.Response().Header().Set("X-Tags", string(js))

	tagSuggest := pub_tagsuggest.TagSuggest(tags)

	return templates.AssertRender(c, http.StatusOK, tagSuggest)
}
