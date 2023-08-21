package admin

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/pkg/admin/dto"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/utils/constants"
	"github.com/muhwyndhamhp/marknotes/utils/markd"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
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

	g.GET("/posts", fe.PostsGet)
	g.GET("/posts_index", fe.PostsIndex)
	g.GET("/posts/new", fe.PostsNew)
	g.POST("/posts/create", fe.PostCreate, htmxMid)
	g.POST("/posts/render", fe.RenderMarkdown, htmxMid)
	g.GET("/posts/:id", fe.GetPostByID)
}

func (fe *PostFrontend) PostsIndex(c echo.Context) error {
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

	posts[len(posts)-1].AppendFormMeta(2)

	return c.Render(http.StatusOK, "posts_index", resp)
}

func (fe *PostFrontend) GetPostByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	post, err := fe.repo.GetByID(ctx, uint(id))
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "posts_detail", post)
}

func (fe *PostFrontend) PostsGet(c echo.Context) error {
	ctx := c.Request().Context()
	page, _ := strconv.Atoi(c.QueryParam(constants.PAGE))
	pageSize, _ := strconv.Atoi(c.QueryParam(constants.PAGE_SIZE))

	posts, err := fe.repo.Get(ctx, scopes.QueryOpts{
		Page:     page,
		PageSize: pageSize,
		Order:    "created_at",
		OrderDir: scopes.Descending,
	})
	if err != nil {
		return err
	}

	if len(posts) > 0 {
		posts[len(posts)-1].AppendFormMeta(page + 1)
	}

	return c.Render(http.StatusOK, "post_list", posts)
}

func (fe *PostFrontend) PostsNew(c echo.Context) error {
	return c.Render(http.StatusOK, "posts_new", nil)
}

func (fe *PostFrontend) RenderMarkdown(c echo.Context) error {
	encoded, err := markd.ParseMD(c.FormValue("content"))
	if err != nil {
		return err
	}
	c.Response().Header().Set("HX-Trigger-After-Swap", "checkTheme")

	return c.HTML(http.StatusOK, encoded)
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

	c.Response().Header().Set("Hx-Redirect", "/")
	return c.JSON(http.StatusOK, nil)
}
