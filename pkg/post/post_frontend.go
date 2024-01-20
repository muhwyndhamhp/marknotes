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
	pub_postlist "github.com/muhwyndhamhp/marknotes/pub/components/postlist"
	pub_post_detail "github.com/muhwyndhamhp/marknotes/pub/pages/post_detail/post_detail"
	pub_post_edit "github.com/muhwyndhamhp/marknotes/pub/pages/post_edit"
	pub_post_index "github.com/muhwyndhamhp/marknotes/pub/pages/post_index"
	pub_post_manage "github.com/muhwyndhamhp/marknotes/pub/pages/post_manage"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	templateRenderer "github.com/muhwyndhamhp/marknotes/template"
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

	search := pub_variables.SearchBar{
		SearchPlaceholder: "Manage Articles...",
		SearchPath:        "/posts?page=1&pageSize=10",
	}

	userID := jwt.AppendAndReturnUserID(c, map[string]interface{}{})

	bodyOpts := pub_variables.BodyOpts{
		HeaderButtons: admin.AppendHeaderButtons(userID),
		FooterButtons: admin.AppendFooterButtons(userID),
		Component:     nil,
	}

	postIndex := pub_post_manage.PostManage(bodyOpts, posts, search)

	return templateRenderer.AssertRender(c, http.StatusOK, postIndex)
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

	if post == nil {
		post = &models.Post{}
	}
	post.FormMeta = map[string]interface{}{
		"SubmitPath": fmt.Sprintf("/posts/%d/update", id),
		"CancelPath": fmt.Sprintf("/posts/%d", id),
	}
	userID := jwt.AppendAndReturnUserID(c, post.FormMeta)
	bodyOpts := pub_variables.BodyOpts{
		HeaderButtons: admin.AppendHeaderButtons(userID),
		FooterButtons: admin.AppendFooterButtons(userID),
		Component:     nil,
	}

	models.SetTagEditable(post.Tags...)

	postEdit := pub_post_edit.PostEdit(bodyOpts, *post)
	return templateRenderer.AssertRender(c, http.StatusOK, postEdit)
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

	search := pub_variables.SearchBar{
		SearchPlaceholder: "Search Articles...",
		SearchPath:        "/posts?page=1&pageSize=10&sortBy=published_at&status=published",
	}

	userID := jwt.AppendAndReturnUserID(c, map[string]interface{}{})

	bodyOpts := pub_variables.BodyOpts{
		HeaderButtons: admin.AppendHeaderButtons(userID),
		FooterButtons: admin.AppendFooterButtons(userID),
		Component:     nil,
	}

	postIndex := pub_post_index.PostIndex(bodyOpts, posts, search)

	return templateRenderer.AssertRender(c, http.StatusOK, postIndex)
}

func (fe *PostFrontend) GetPostBySlug(c echo.Context) error {
	ctx := c.Request().Context()

	slug := strings.TrimSpace(c.Param("slug"))

	posts, err := fe.repo.Get(ctx, scopes.Where("slug = ?", slug))
	if err != nil {
		return err
	}

	post := posts[0]

	userID := jwt.AppendAndReturnUserID(c, map[string]interface{}{})

	post.FormMeta = map[string]interface{}{
		"UserID": userID,
	}
	bodyOpts := pub_variables.BodyOpts{
		HeaderButtons: admin.AppendHeaderButtons(userID),
		FooterButtons: admin.AppendFooterButtons(userID),
		Component:     nil,
	}

	postDetail := pub_post_detail.PostDetail(bodyOpts, post)

	return templateRenderer.AssertRender(c, http.StatusOK, postDetail)
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
	var bodyOpts pub_variables.BodyOpts

	if post.FormMeta == nil {
		post.FormMeta = map[string]interface{}{}
	}

	if claims != nil {
		post.FormMeta["UserID"] = claims.UserID
		bodyOpts = pub_variables.BodyOpts{
			HeaderButtons: admin.AppendHeaderButtons(claims.UserID),
			FooterButtons: admin.AppendFooterButtons(claims.UserID),
			Component:     nil,
		}
	} else {
		post.FormMeta["UserID"] = claims.UserID
		bodyOpts = pub_variables.BodyOpts{
			HeaderButtons: admin.AppendHeaderButtons(claims.UserID),
			FooterButtons: admin.AppendFooterButtons(claims.UserID),
			Component:     nil,
		}
	}

	post.FormMeta["CenterAlign"] = true

	if post.Status == values.Draft &&
		(claims == nil || claims.UserID == 0) {
		return c.Redirect(http.StatusFound, "/")
	}

	postDetail := pub_post_detail.PostDetail(bodyOpts, *post)

	return templateRenderer.AssertRender(c, http.StatusOK, postDetail)
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

	postList := pub_postlist.PostList(posts)

	return templateRenderer.AssertRender(c, http.StatusOK, postList)
}

func (fe *PostFrontend) PostsNew(c echo.Context) error {
	post := &models.Post{}

	post.FormMeta = map[string]interface{}{
		"SubmitPath": "/posts/create",
		"CancelPath": "/posts_manage",
	}

	userID := jwt.AppendAndReturnUserID(c, post.FormMeta)
	bodyOpts := pub_variables.BodyOpts{
		HeaderButtons: admin.AppendHeaderButtons(userID),
		FooterButtons: admin.AppendFooterButtons(userID),
		Component:     nil,
	}

	models.SetTagEditable(post.Tags...)

	postEdit := pub_post_edit.PostEdit(bodyOpts, *post)
	return templateRenderer.AssertRender(c, http.StatusOK, postEdit)
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
