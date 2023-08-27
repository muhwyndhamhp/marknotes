package post

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/middlewares"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/pkg/post/dto"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	"github.com/muhwyndhamhp/marknotes/utils/jwt"
	"github.com/muhwyndhamhp/marknotes/utils/markd"
	"github.com/muhwyndhamhp/marknotes/utils/params"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"gorm.io/gorm"
)

type PostFrontend struct {
	repo models.PostRepository
}

func NewPostFrontend(g *echo.Group,
	repo models.PostRepository,
	htmxMid echo.MiddlewareFunc,
	authMid echo.MiddlewareFunc,
	authDescribeMid echo.MiddlewareFunc,
	byIDMiddleware echo.MiddlewareFunc,
) {
	fe := &PostFrontend{
		repo: repo,
	}

	g.GET("/posts", fe.PostsGet, authDescribeMid)
	g.GET("/posts_index", fe.PostsIndex, authDescribeMid)
	g.GET("/posts_manage", fe.PostsManage, authMid)

	g.GET("/posts/new", fe.PostsNew, authMid)
	g.POST("/posts/create", fe.PostCreate, htmxMid, authMid)
	g.POST("/posts/render", fe.RenderMarkdown, htmxMid, authMid)

	g.GET("/posts/:id", fe.GetPostByID, authDescribeMid, byIDMiddleware)
	g.GET("/posts/:id/edit", fe.PostEdit, authMid, byIDMiddleware)
	g.POST("/posts/:id/update", fe.PostUpdate, htmxMid, authMid, byIDMiddleware)
	g.GET("/posts/:id/delete", fe.PostDelete, htmxMid, authMid, byIDMiddleware)
	g.GET("/posts/:id/publish", fe.PostPublish, htmxMid, authMid, byIDMiddleware)
	g.GET("/posts/:id/draft", fe.PostDraft, htmxMid, byIDMiddleware)
}

func (fe *PostFrontend) PostDraft(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := c.Get(middlewares.ByIDKey).(int)

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
		scopes.OrderBy("created_at", scopes.Descending),
	)
	if err != nil {
		return err
	}
	if len(posts) > 0 {
		posts[len(posts)-1].AppendFormMeta(2, false, "")
	}

	resp := map[string]interface{}{"Posts": posts}
	jwt.AppendUserID(c, resp)

	return c.Render(http.StatusOK, "posts_index", resp)
}

func (fe *PostFrontend) PostPublish(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := c.Get(middlewares.ByIDKey).(int)

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
	id, _ := c.Get(middlewares.ByIDKey).(int)

	if err := fe.repo.Delete(ctx, uint(id)); err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Redirect", "/")
	return c.JSON(http.StatusOK, nil)
}

func (fe *PostFrontend) PostUpdate(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := c.Get(middlewares.ByIDKey).(int)

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
	post.Tags = []*models.Tag{}

	for i := range req.Tags {
		id, _ := strconv.Atoi(req.Tags[i])
		post.Tags = append(post.Tags, &models.Tag{
			Model: gorm.Model{
				ID: uint(id),
			},
		})
	}

	if err = fe.repo.Upsert(ctx, post); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/posts/%d", post.ID))
}

func (fe *PostFrontend) PostEdit(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := c.Get(middlewares.ByIDKey).(int)

	post, err := fe.repo.GetByID(ctx, uint(id))
	if err != nil {
		return err
	}

	post.FormMeta = map[string]interface{}{}
	jwt.AppendUserID(c, post.FormMeta)

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
	if len(posts) > 0 {
		posts[len(posts)-1].AppendFormMeta(2, true, "published_at")
	}

	resp := map[string]interface{}{"Posts": posts}
	jwt.AppendUserID(c, resp)

	return c.Render(http.StatusOK, "posts_index", resp)
}

func (fe *PostFrontend) GetPostByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := c.Get(middlewares.ByIDKey).(int)

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

	if post.Status == values.Draft &&
		(claims == nil || claims.UserID == 0) {
		return c.Redirect(http.StatusFound, "/")
	}

	return c.Render(http.StatusOK, "posts_detail", post)
}

func (fe *PostFrontend) PostsGet(c echo.Context) error {
	ctx := c.Request().Context()
	page, pageSize, sortBy, statusStr := params.GetCommonParams(c)
	status := values.PostStatus(statusStr)
	if sortBy == "" {
		sortBy = "created_at"
	}

	posts, err := fe.repo.Get(ctx,
		scopes.Paginate(page, pageSize),
		scopes.WithStatus(status),
		scopes.OrderBy(sortBy, scopes.Descending))
	if err != nil {
		return err
	}

	onlyPublised := status == values.Published
	if len(posts) > 0 {
		posts[len(posts)-1].AppendFormMeta(page+1, onlyPublised, sortBy)
	}

	return c.Render(http.StatusOK, "post_list", posts)
}

func (fe *PostFrontend) PostsNew(c echo.Context) error {
	post := &models.Post{}

	post.FormMeta = map[string]interface{}{}
	jwt.AppendUserID(c, post.FormMeta)

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

	claims, _ := c.Get(jwt.AuthClaimKey).(*jwt.Claims)
	if claims == nil {
		return c.JSON(http.StatusBadRequest, errors.New("claim is empty"))
	}

	post := models.Post{
		Title:          req.Title,
		Content:        content,
		EncodedContent: template.HTML(encoded),
		Status:         values.Draft,
		UserID:         claims.UserID,
	}

	for i := range req.Tags {
		id, _ := strconv.Atoi(req.Tags[i])
		post.Tags = append(post.Tags, &models.Tag{
			Model: gorm.Model{
				ID: uint(id),
			},
		})
	}

	err = fe.repo.Upsert(ctx, &post)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Redirect", fmt.Sprintf("/posts/%d", post.ID))
	return c.JSON(http.StatusOK, nil)
}
