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
}

func (fe *PostFrontend) PostsNew(c echo.Context) error {
	return c.Render(http.StatusOK, "posts_new", nil)
}
