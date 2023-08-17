package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/gotes-mx/pkg/models"
)

type PostFrontend struct {
	repo models.PostRepository
}

func NewPostFrontend(g *echo.Group, repo models.PostRepository) {
	fe := &PostFrontend{
		repo: repo,
	}

	g.GET("/posts/new", fe.PostsNew)
	g.GET("/posts/new/cancel", fe.PostNewCancel)
}

func (fe *PostFrontend) PostsNew(c echo.Context) error {
	return c.Render(http.StatusOK, "posts_new", nil)
}

func (fe *PostFrontend) PostNewCancel(c echo.Context) error {
	c.Response().Header().Set("HX-Redirect", "/admin")
	return c.Render(http.StatusOK, "admin_index", nil)
}
