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
	"github.com/muhwyndhamhp/marknotes/pkg/admin"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/pkg/post/dto"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	"github.com/muhwyndhamhp/marknotes/utils/jwt"
	"github.com/muhwyndhamhp/marknotes/utils/markd"
	"github.com/muhwyndhamhp/marknotes/utils/params"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"github.com/muhwyndhamhp/marknotes/utils/strman"
	"github.com/muhwyndhamhp/marknotes/utils/tern"
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

	// Alias `articles` for `posts`
	g.GET("/articles", fe.PostsIndex, authDescribeMid)
	g.GET("/articles/:slug", fe.GetPostBySlug, authDescribeMid)
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
		scopes.PostIndexedOnly(),
	)
	if err != nil {
		return err
	}
	if len(posts) > 0 {
		posts[len(posts)-1].AppendFormMeta(2, values.None, "", "")
	}

	resp := map[string]interface{}{"Posts": posts}
	userID := jwt.AppendAndReturnUserID(c, resp)
	resp[admin.HeaderButtonsKey] = admin.AppendHeaderButtons(userID)
	resp[admin.FooterButtonsKey] = admin.AppendFooterButtons(userID)
	resp[SearchBarKey] = SearchBar{
		SearchPlaceholder: "Manage Articles...",
		SearchPath:        "/posts?page=1&pageSize=10",
	}

	return c.Render(http.StatusOK, "posts_manage", resp)
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

	if post.Slug == "" {
		slug, err := strman.GenerateSlug(post.Title)
		if err != nil {
			return err
		}
		post.Slug = slug
	}

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

	post.FormMeta = map[string]interface{}{
		"SubmitPath": fmt.Sprintf("/posts/%d/update", id),
		"CancelPath": fmt.Sprintf("/posts/%d", id),
	}
	userID := jwt.AppendAndReturnUserID(c, post.FormMeta)
	post.FormMeta[admin.HeaderButtonsKey] = admin.AppendHeaderButtons(userID)
	post.FormMeta[admin.FooterButtonsKey] = admin.AppendFooterButtons(userID)

	models.SetTagEditable(post.Tags...)

	return c.Render(http.StatusOK, "posts_edit", post)
}

func (fe *PostFrontend) PostsIndex(c echo.Context) error {
	ctx := c.Request().Context()

	posts, err := fe.repo.Get(ctx,
		scopes.Paginate(1, 10),
		scopes.OrderBy("published_at", scopes.Descending),
		scopes.WithStatus(values.Published),
		scopes.PostIndexedOnly(),
	)
	if err != nil {
		return err
	}
	if len(posts) > 0 {
		posts[len(posts)-1].AppendFormMeta(2, values.Published, "published_at", "")
	}

	resp := map[string]interface{}{"Posts": posts}
	userID := jwt.AppendAndReturnUserID(c, resp)
	resp[admin.HeaderButtonsKey] = admin.AppendHeaderButtons(userID)
	resp[admin.FooterButtonsKey] = admin.AppendFooterButtons(userID)
	resp[SearchBarKey] = SearchBar{
		SearchPlaceholder: "Search Articles...",
		SearchPath:        "/posts?page=1&pageSize=10&sortBy=published_at&status=published",
	}

	return c.Render(http.StatusOK, "posts_index", resp)
}

func (fe *PostFrontend) GetPostBySlug(c echo.Context) error {
	ctx := c.Request().Context()

	slug := strings.TrimSpace(c.Param("slug"))

	post, err := fe.repo.Get(ctx, scopes.Where("slug = ?", slug))
	if err != nil {
		return err
	}

	return fe.renderPost(c, &post[0])
}

func (fe *PostFrontend) GetPostByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := c.Get(middlewares.ByIDKey).(int)

	post, err := fe.repo.GetByID(ctx, uint(id))
	if err != nil {
		return err
	}

	if post.Slug == "" {
		return fe.renderPost(c, post)
	} else {
		return c.Redirect(http.StatusFound, fmt.Sprintf("/articles/%s", post.Slug))
	}
}

func (fe *PostFrontend) renderPost(c echo.Context, post *models.Post) error {
	claims, _ := c.Get(jwt.AuthClaimKey).(*jwt.Claims)
	if claims != nil {
		post.FormMeta = map[string]interface{}{
			"UserID":               claims.UserID,
			admin.HeaderButtonsKey: admin.AppendHeaderButtons(claims.UserID),
			admin.FooterButtonsKey: admin.AppendFooterButtons(claims.UserID),
		}
	} else {
		post.FormMeta = map[string]interface{}{
			admin.HeaderButtonsKey: admin.AppendHeaderButtons(0),
			admin.FooterButtonsKey: admin.AppendFooterButtons(0),
		}
	}

	post.FormMeta["CenterAlign"] = true

	if post.Status == values.Draft &&
		(claims == nil || claims.UserID == 0) {
		return c.Redirect(http.StatusFound, "/")
	}

	return c.Render(http.StatusOK, "posts_detail", post)
}

func (fe *PostFrontend) PostsGet(c echo.Context) error {
	ctx := c.Request().Context()
	page, pageSize, sortBy, statusStr, keyword, loadNext := params.GetCommonParams(c)
	status := values.PostStatus(statusStr)

	scp := []scopes.QueryScope{
		scopes.OrderBy(tern.String(sortBy, "created_at"), scopes.Descending),
		scopes.Paginate(page, pageSize),
	}

	if keyword != "" {
		scp = append(scp, scopes.PostDeepMatch(keyword, status))
	} else {
		scp = append(scp, scopes.WithStatus(status))
	}

	posts, err := fe.repo.Get(ctx, scp...)
	if err != nil {
		return err
	}

	if len(posts) > 0 && loadNext {
		posts[len(posts)-1].AppendFormMeta(page+1, status, sortBy, keyword)
	}

	return c.Render(http.StatusOK, "post_list", posts)
}

func (fe *PostFrontend) PostsNew(c echo.Context) error {
	post := &models.Post{}

	post.FormMeta = map[string]interface{}{
		"SubmitPath": "/posts/create",
		"CancelPath": "/posts_manage",
	}
	userID := jwt.AppendAndReturnUserID(c, post.FormMeta)
	post.FormMeta[admin.HeaderButtonsKey] = admin.AppendHeaderButtons(userID)
	post.FormMeta[admin.FooterButtonsKey] = admin.AppendFooterButtons(userID)

	models.SetTagEditable(post.Tags...)

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

	slug, err := strman.GenerateSlug(req.Title)
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
		Slug:           slug,
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
