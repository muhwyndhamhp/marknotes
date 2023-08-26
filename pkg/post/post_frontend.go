package post

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/pkg/admin/dto"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	"github.com/muhwyndhamhp/marknotes/utils/constants"
	"github.com/muhwyndhamhp/marknotes/utils/jwt"
	"github.com/muhwyndhamhp/marknotes/utils/markd"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
)

type PostFrontend struct {
	repo    models.PostRepository
	htmxMid echo.MiddlewareFunc
}

func NewPostFrontend(g *echo.Group,
	repo models.PostRepository,
	htmxMid echo.MiddlewareFunc,
	authMid echo.MiddlewareFunc,
	authDescribeMid echo.MiddlewareFunc) {
	fe := &PostFrontend{
		repo:    repo,
		htmxMid: htmxMid,
	}

	g.GET("/posts", fe.PostsGet, authDescribeMid)
	g.GET("/posts_index", fe.PostsIndex, authDescribeMid)
	g.GET("/posts_manage", fe.PostsManage, authMid)

	g.GET("/posts/new", fe.PostsNew, authMid)
	g.POST("/posts/create", fe.PostCreate, htmxMid, authMid)
	g.POST("/posts/render", fe.RenderMarkdown, htmxMid, authMid)

	g.GET("/posts/:id", fe.GetPostByID, authDescribeMid)
	g.GET("/posts/:id/edit", fe.PostEdit, authMid)
	g.POST("/posts/:id/update", fe.PostUpdate, htmxMid, authMid)
	g.GET("/posts/:id/delete", fe.PostDelete, htmxMid, authMid)
	g.GET("/posts/:id/publish", fe.PostPublish, htmxMid, authMid)
	g.GET("/posts/:id/draft", fe.PostDraft, htmxMid)
}

func (fe *PostFrontend) PostDraft(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	post, err := fe.repo.GetByID(ctx, uint(id))
	if err != nil {
		return err
	}

	post.Status = values.Draft

	if err = fe.repo.Upsert(ctx, post); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/posts/%d", post.ID))
}

func (fe *PostFrontend) PostsManage(c echo.Context) error {
	ctx := c.Request().Context()
	posts, err := fe.repo.Get(ctx,
		scopes.Paginate(1, 10),
		scopes.OrderBy("published_at", scopes.Descending),
	)

	if err != nil {
		return err
	}
	resp := map[string]interface{}{
		"Posts": posts,
	}

	claims, _ := c.Get(jwt.AuthClaimKey).(*jwt.Claims)

	if claims != nil {
		resp["UserID"] = claims.UserID
	}

	posts[len(posts)-1].AppendFormMeta(2, false)

	return c.Render(http.StatusOK, "posts_index", resp)
}
func (fe *PostFrontend) PostPublish(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	post, err := fe.repo.GetByID(ctx, uint(id))
	if err != nil {
		return err
	}

	post.Status = values.Published
	post.PublishedAt = time.Now()

	if err = fe.repo.Upsert(ctx, post); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/posts/%d", post.ID))
}

func (fe *PostFrontend) PostDelete(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := fe.repo.Delete(ctx, uint(id)); err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Redirect", "/")
	return c.JSON(http.StatusOK, nil)
}

func (fe *PostFrontend) PostUpdate(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	var req dto.PostCreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	content := strings.TrimSpace(req.Content)

	encoded, err := markd.ParseMD(content)
	if err != nil {
		return err
	}

	post, err := fe.repo.GetByID(ctx, uint(id))
	if err != nil {
		return err
	}

	post.Title = req.Title
	post.Content = content
	post.EncodedContent = template.HTML(encoded)

	if err = fe.repo.Upsert(ctx, post); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/posts/%d", post.ID))
}

func (fe *PostFrontend) PostEdit(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	post, err := fe.repo.GetByID(ctx, uint(id))
	if err != nil {
		return err
	}

	claims, _ := c.Get(jwt.AuthClaimKey).(*jwt.Claims)

	if claims != nil {
		post.FormMeta = map[string]interface{}{
			"UserID": claims.UserID,
		}
	}

	return c.Render(http.StatusOK, "posts_edit", post)
}

func (fe *PostFrontend) PostsIndex(c echo.Context) error {
	ctx := c.Request().Context()

	posts, err := fe.repo.Get(ctx,
		scopes.Paginate(1, 10),
		scopes.OrderBy("published_at", scopes.Descending),
		scopes.WithStatus(values.Published),
	)

	if err != nil {
		return err
	}
	resp := map[string]interface{}{
		"Posts": posts,
	}

	claims, _ := c.Get(jwt.AuthClaimKey).(*jwt.Claims)

	if claims != nil {
		resp["UserID"] = claims.UserID
	}

	posts[len(posts)-1].AppendFormMeta(2, true)

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
	claims, _ := c.Get(jwt.AuthClaimKey).(*jwt.Claims)

	if claims != nil {
		post.FormMeta = map[string]interface{}{
			"UserID": claims.UserID,
		}
	}

	return c.Render(http.StatusOK, "posts_detail", post)
}

func (fe *PostFrontend) PostsGet(c echo.Context) error {
	ctx := c.Request().Context()
	page, _ := strconv.Atoi(c.QueryParam(constants.PAGE))
	pageSize, _ := strconv.Atoi(c.QueryParam(constants.PAGE_SIZE))
	statusStr := c.QueryParam(constants.STATUS)
	status := values.PostStatus(statusStr)

	posts, err := fe.repo.Get(ctx,
		scopes.Paginate(page, pageSize),
		scopes.OrderBy("published_at", scopes.Descending),
		scopes.WithStatus(status),
	)
	if err != nil {
		return err
	}

	onlyPublised := status == values.Published
	if len(posts) > 0 {
		posts[len(posts)-1].AppendFormMeta(page+1, onlyPublised)
	}

	return c.Render(http.StatusOK, "post_list", posts)
}

func (fe *PostFrontend) PostsNew(c echo.Context) error {
	post := &models.Post{}

	claims, _ := c.Get(jwt.AuthClaimKey).(*jwt.Claims)

	if claims != nil {
		post.FormMeta = map[string]interface{}{
			"UserID": claims.UserID,
		}
	}
	return c.Render(http.StatusOK, "posts_new", post)
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

	content := strings.TrimSpace(req.Content)

	encoded, err := markd.ParseMD(content)
	if err != nil {
		return err
	}

	post := models.Post{
		Title:          req.Title,
		Content:        content,
		EncodedContent: template.HTML(encoded),
		Status:         values.Draft,
	}

	err = fe.repo.Upsert(ctx, &post)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Redirect", "/")
	return c.JSON(http.StatusOK, nil)
}