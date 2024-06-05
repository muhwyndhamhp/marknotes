package articles

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	templates "github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/constants"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"github.com/muhwyndhamhp/marknotes/utils/tern"
	"gorm.io/gorm"
)

type Handler struct{}

func NewHandler(g echo.Group, db *gorm.DB) {
}

func (h *Handler) Index(c echo.Context) error {
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
	if len(posts) <= 0 && page > 1 {
		appendRoute := ""
		if source == constants.TARGET_SOURCE_PARTIAL {
			appendRoute = "&source=source-partial"
		}
		path := fmt.Sprintf("/dashboard/articles?page=%d&pageSize=%d%s", page-1, pageSize, appendRoute)
		fmt.Println(path)
		return c.Redirect(http.StatusFound, path)
	}

	opts := pub_variables.DashboardOpts{
		Nav:         nav(0),
		BreadCrumbs: fe.Breadcrumbs("dashboard/articles"),
	}

	pageSizes := fe.SizeDropdown(page, pageSize)
	pages := fe.PageDropdown(tern.Int(page, 1), tern.Int(pageSize, 10), count)
	articleVM := ArticlesVM{
		Opts:       opts,
		Posts:      posts,
		PageSizes:  pageSizes,
		Pages:      pages,
		CreatePath: "/dashboard/articles/new",
	}

	dashboard := Articles(articleVM)

	if !partial {
		return templates.AssertRender(c, http.StatusOK, dashboard)
	} else {
		articles := ArticleOOB(posts, pageSizes, pages)
		return templates.AssertRender(c, http.StatusOK, articles)
	}
}
