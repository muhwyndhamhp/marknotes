package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
)

type AdminFrontend struct {
	repo models.PostRepository
}

func NewAdminFrontend(g *echo.Group, repo models.PostRepository) {
	fe := &AdminFrontend{
		repo: repo,
	}

	g.GET("", fe.Index)
}

func (fe *AdminFrontend) Index(c echo.Context) error {
	ctx := c.Request().Context()

	posts, err := fe.repo.Get(ctx, scopes.QueryOpts{
		Page:     1,
		PageSize: 10,
		Order:    "created_at",
		OrderDir: scopes.Descending,
	})

	if err != nil {
		return err
	}
	resp := map[string]interface{}{
		"Posts": posts,
	}

	return c.Render(http.StatusOK, "admin_index", resp)
}
