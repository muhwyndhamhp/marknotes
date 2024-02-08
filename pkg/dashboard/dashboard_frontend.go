package dashboard

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	pub_dashboards_articles "github.com/muhwyndhamhp/marknotes/pub/pages/dashboards/articles"
	pub_dashboard_editor "github.com/muhwyndhamhp/marknotes/pub/pages/dashboards/editor"
	pub_dashboards_profile "github.com/muhwyndhamhp/marknotes/pub/pages/dashboards/profile"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	"github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/constants"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"github.com/muhwyndhamhp/marknotes/utils/tern"
)

type DashboardFrontend struct {
	PostRepo models.PostRepository
}

func NewDashboardFrontend(
	g *echo.Group,
	PostRepo models.PostRepository,
	htmxMid echo.MiddlewareFunc,
	authMid echo.MiddlewareFunc,
	authDescribeMid echo.MiddlewareFunc,
	byIDMiddleware echo.MiddlewareFunc,
) {
	fe := &DashboardFrontend{PostRepo}

	g.GET("/dashboard/articles", fe.Articles, authMid)
	g.GET("/dashboard/profile", fe.Profile, authMid)
	g.GET("/dashboard/editor", fe.Editor, authMid)
}

func (fe *DashboardFrontend) Editor(c echo.Context) error {

	opts := pub_variables.DashboardOpts{Nav: nav(1)}

	dashboard := pub_dashboard_editor.Editor(opts)

	return template.AssertRender(c, http.StatusOK, dashboard)
}

func (fe *DashboardFrontend) Profile(c echo.Context) error {
	opts := pub_variables.DashboardOpts{Nav: nav(1)}

	dashboard := pub_dashboards_profile.Profile(opts)

	return template.AssertRender(c, http.StatusOK, dashboard)
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
		Opts:      opts,
		Posts:     posts,
		PageSizes: pageSizes,
		Pages:     pages,
	}

	dashboard := pub_dashboards_articles.Articles(articleVM)

	if !partial {
		return template.AssertRender(c, http.StatusOK, dashboard)
	} else {
		articles := pub_dashboards_articles.ArticleOOB(posts, pageSizes, pages)
		return template.AssertRender(c, http.StatusOK, articles)
	}
}

func nav(indexSelected int) []pub_variables.DrawerMenu {
	lists := []pub_variables.DrawerMenu{
		{
			Label:    "Articles",
			URL:      "/dashboard/articles",
			IsActive: false,
		},
		{
			Label:    "Profile",
			URL:      "/dashboard/profile",
			IsActive: false,
		},
		{
			Label:    "Create Post",
			URL:      "/dashboard/editor",
			IsActive: false,
		},
	}

	lists[indexSelected].IsActive = true

	return lists
}
