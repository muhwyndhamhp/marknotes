package admin

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/pkg/admin/dto"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/utils/markd"
)

type PostFrontend struct {
	repo    models.PostRepository
	htmxMid echo.MiddlewareFunc
}

func NewPostFrontend(g *echo.Group, repo models.PostRepository, htmxMid echo.MiddlewareFunc) {
	fe := &PostFrontend{
		repo:    repo,
		htmxMid: htmxMid,
	}

	g.GET("/posts/new", fe.PostsNew, htmxMid)
	g.POST("/posts/create", fe.PostCreate, htmxMid)
}

func (fe *PostFrontend) PostsNew(c echo.Context) error {
	return c.Render(http.StatusOK, "posts_new", nil)
}

func (fe *PostFrontend) PostCreate(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.PostCreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	encoded, err := markd.ParseMD(req.Content)
	if err != nil {
		return err
	}

	post := models.Post{
		Title:          req.Title,
		Content:        req.Content,
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
