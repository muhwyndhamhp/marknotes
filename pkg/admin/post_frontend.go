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
	g.POST("/posts/create", fe.PostCreate)
}

func (fe *PostFrontend) PostsNew(c echo.Context) error {
	return c.Render(http.StatusOK, "posts_new", nil)
}

func (fe *PostFrontend) PostCreate(c echo.Context) error {
	ctx := c.Request().Context()
	post := models.Post{
		Title:          c.FormValue("title"),
		Content:        c.FormValue("content"),
		EncodedContent: "",
		Status:         models.Draft,
	}

	err := fe.repo.Upsert(ctx, &post)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Redirect", "/admin")
	return c.JSON(http.StatusOK, nil)
}
