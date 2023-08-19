package admin

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/gotes-mx/pkg/models"
	"github.com/muhwyndhamhp/gotes-mx/utils/markd"
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

	encoded, err := markd.ParseMD(c.FormValue("content"))
	if err != nil {
		return err
	}

	post := models.Post{
		Title:          c.FormValue("title"),
		Content:        encoded,
		EncodedContent: template.HTML(encoded),
		Status:         models.Draft,
	}

	err = fe.repo.Upsert(ctx, &post)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Redirect", "/admin")
	return c.JSON(http.StatusOK, nil)
}
