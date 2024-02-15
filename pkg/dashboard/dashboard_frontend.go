package dashboard

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	pub_editor "github.com/muhwyndhamhp/marknotes/pub/components/editor"
	pub_tagsuggest "github.com/muhwyndhamhp/marknotes/pub/components/tagsuggest"
	pub_dashboard "github.com/muhwyndhamhp/marknotes/pub/pages/dashboards"
	pub_dashboards_articles "github.com/muhwyndhamhp/marknotes/pub/pages/dashboards/articles"
	pub_dashboards_articles_new "github.com/muhwyndhamhp/marknotes/pub/pages/dashboards/articles/create"
	pub_dashboards_profile "github.com/muhwyndhamhp/marknotes/pub/pages/dashboards/profile"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	templates "github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/constants"
	"github.com/muhwyndhamhp/marknotes/utils/jwt"
	"github.com/muhwyndhamhp/marknotes/utils/sanitizations"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"github.com/muhwyndhamhp/marknotes/utils/strman"
	"github.com/muhwyndhamhp/marknotes/utils/tern"
	"gorm.io/gorm"
)

type DashboardFrontend struct {
	PostRepo models.PostRepository
	TagRepo  models.TagRepository
}

func NewDashboardFrontend(
	g *echo.Group,
	PostRepo models.PostRepository,
	TagRepo models.TagRepository,
	htmxMid echo.MiddlewareFunc,
	authMid echo.MiddlewareFunc,
	authDescribeMid echo.MiddlewareFunc,
	byIDMiddleware echo.MiddlewareFunc,
) {
	fe := &DashboardFrontend{PostRepo, TagRepo}

	g.GET("/dashboard/articles", fe.Articles, authMid)
	g.POST("/dashboard/articles/create", fe.Create)
	g.GET("/dashboard/articles/new", fe.ArticlesCreate, authMid)
	g.GET("/dashboard/profile", fe.Profile, authMid)
	g.GET("/dashboard/editor", fe.Editor, authMid)
	g.GET("/dashboard/tags", fe.Tags, authMid)
}

type ArticlesCreateRequest struct {
	Title   string   `json:"title" validate:"required" form:"title"`
	Content string   `json:"content" validate:"required" form:"content"`
	Tags    []string `json:"tags" form:"tags"`
}

func (fe *DashboardFrontend) Tags(c echo.Context) error {
	ctx := c.Request().Context()

	tagQuery := c.QueryParam("tag")

	tagName := strings.ToLower(strings.TrimSpace(tagQuery))
	tagSlug := strings.ReplaceAll(tagName, " ", "-")

	tags, err := fe.TagRepo.Get(
		ctx,
		scopes.Where("slug LIKE ?", fmt.Sprintf("%%%s%%", tagSlug)),
		scopes.Paginate(1, 5),
	)
	if err != nil {
		return err
	}

	tagSuggest := pub_tagsuggest.TagSuggest(tags)

	return templates.AssertRender(c, http.StatusOK, tagSuggest)
}

func (fe *DashboardFrontend) Create(c echo.Context) error {
	// ctx := c.Request().Context()

	var req ArticlesCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	sanitizedHTMl := sanitizations.SanitizeHtml(req.Content)

	slug, err := strman.GenerateSlug(req.Title)
	if err != nil {
		return err
	}

	claims, _ := c.Get(jwt.AuthClaimKey).(*jwt.Claims)
	if claims == nil {
		return c.JSON(http.StatusBadRequest, errors.New("claim is empty"))
	}

	render := pub_dashboard.ExampleRaw(sanitizedHTMl)

	post := models.Post{
		Title:          req.Title,
		EncodedContent: template.HTML(sanitizedHTMl),
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

	return templates.AssertRenderLog(c, http.StatusOK, render)

	// return c.HTML(http.StatusOK, sanitizedHTMl)
}

func (fe *DashboardFrontend) ArticlesCreate(c echo.Context) error {
	opts := pub_variables.DashboardOpts{
		Nav: nav(2),
	}

	baseURL := strings.Split(config.Get(config.OAUTH_URL), "/callback")[0]
	uploadURL := fmt.Sprintf("%s/posts/%d/media/upload", baseURL, 0)

	vm := pub_dashboards_articles_new.NewVM{
		Opts:      opts,
		UploadURL: uploadURL,
	}
	articlesNew := pub_dashboards_articles_new.New(vm)

	return templates.AssertRender(c, http.StatusOK, articlesNew)
}

func (fe *DashboardFrontend) Editor(c echo.Context) error {
	baseURL := strings.Split(config.Get(config.OAUTH_URL), "/callback")[0]
	uploadURL := fmt.Sprintf("%s/posts/%d/media/upload", baseURL, 0)

	dashboard := pub_editor.Editor(uploadURL)

	return templates.AssertRender(c, http.StatusOK, dashboard)
}

func (fe *DashboardFrontend) Profile(c echo.Context) error {
	opts := pub_variables.DashboardOpts{Nav: nav(1)}

	dashboard := pub_dashboards_profile.Profile(opts)

	return templates.AssertRender(c, http.StatusOK, dashboard)
}

func (fe *DashboardFrontend) SizeDropdown(page, pageSize int) pub_variables.DropdownVM {
	arrays := []pub_variables.DropdownItem{}
	item := 0
	for i := range []int{0, 1, 2} {
		size := (i + 1) * 10
		arrays = append(arrays, pub_variables.DropdownItem{
			Label:  fmt.Sprintf("%d", size),
			Path:   fmt.Sprintf("/dashboard/articles?page=%d&pageSize=%d&source=source-partial", page, size),
			Target: "#articles",
		})
		if size == pageSize {
			item = i
		}
	}
	return pub_variables.DropdownVM{
		Items:    arrays,
		Selected: item,
	}
}

func (fe *DashboardFrontend) PageDropdown(page, pageSize, count int) pub_variables.DropdownVM {
	arrays := []pub_variables.DropdownItem{}
	item := 0
	for i := 0; (i)*pageSize <= count; i++ {
		currentPage := i + 1
		size := pageSize
		arrays = append(arrays, pub_variables.DropdownItem{
			Label:  fmt.Sprintf("%d", currentPage),
			Path:   fmt.Sprintf("/dashboard/articles?page=%d&pageSize=%d&source=source-partial", currentPage, size),
			Target: "#articles",
		})
		if currentPage == page {
			item = i
		}
	}

	return pub_variables.DropdownVM{
		Items:    arrays,
		Selected: item,
	}
}

func (fe *DashboardFrontend) Articles(c echo.Context) error {
	ctx := c.Request().Context()
	page, _ := strconv.Atoi(c.QueryParam(constants.PAGE))
	pageSize, _ := strconv.Atoi(c.QueryParam(constants.PAGE_SIZE))
	source := c.QueryParam(constants.TARGET_SOURCE)

	hx_request, _ := strconv.ParseBool(c.Request().Header.Get("Hx-Request"))

	partial := source == constants.TARGET_SOURCE_PARTIAL && hx_request

	count := fe.PostRepo.Count(ctx)

	posts, err := fe.PostRepo.Get(ctx,
		scopes.Paginate(page, pageSize),
		scopes.OrderBy("created_at", scopes.Descending),
		scopes.PostIndexedOnly(),
	)
	if err != nil {
		return err
	}
	if len(posts) > 0 {
		posts[len(posts)-1].AppendFormMeta(2, values.None, "", "")
	}

	opts := pub_variables.DashboardOpts{Nav: nav(0)}

	pageSizes := fe.SizeDropdown(page, pageSize)
	pages := fe.PageDropdown(tern.Int(page, 1), tern.Int(pageSize, 10), count)
	articleVM := pub_dashboards_articles.ArticlesVM{
		Opts:       opts,
		Posts:      posts,
		PageSizes:  pageSizes,
		Pages:      pages,
		CreatePath: "/dashboard/articles/new",
	}

	dashboard := pub_dashboards_articles.Articles(articleVM)

	if !partial {
		return templates.AssertRender(c, http.StatusOK, dashboard)
	} else {
		articles := pub_dashboards_articles.ArticleOOB(posts, pageSizes, pages)
		return templates.AssertRender(c, http.StatusOK, articles)
	}
}

func nav(indexSelected int) []pub_variables.DrawerMenu {
	lists := []pub_variables.DrawerMenu{
		{
			Label:     "Articles",
			URL:       "/dashboard/articles",
			IsActive:  false,
			IsBoosted: true,
		},
		{
			Label:     "Profile",
			URL:       "/dashboard/profile",
			IsActive:  false,
			IsBoosted: true,
		},
		{
			Label:     "Create Post",
			URL:       "/dashboard/editor",
			IsActive:  false,
			IsBoosted: false,
		},
	}

	lists[indexSelected].IsActive = true

	return lists
}
