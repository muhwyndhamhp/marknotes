package articles

import (
	"fmt"
	"github.com/muhwyndhamhp/marknotes/web/routes/dashboard/_partials/breadcrumb"
	"github.com/muhwyndhamhp/marknotes/web/routes/dashboard/_partials/sidebar"
	"github.com/muhwyndhamhp/marknotes/web/routes/dashboard/articles/_partials/page"
	"github.com/muhwyndhamhp/marknotes/web/routes/dashboard/articles/_partials/size"
	new2 "github.com/muhwyndhamhp/marknotes/web/routes/dashboard/articles/new"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/db"
	"github.com/muhwyndhamhp/marknotes/models"
	"github.com/muhwyndhamhp/marknotes/models/post"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	templates "github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/constants"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"github.com/muhwyndhamhp/marknotes/utils/tern"
	"gorm.io/gorm"
)

type Handler struct {
	DB         *gorm.DB
	NewHandler *new2.Handler
}

func NewHandler(
	g *echo.Group,
	db *gorm.DB,
	authMid echo.MiddlewareFunc,
	authDescribeMid echo.MiddlewareFunc,
) *Handler {
	path := "/articles"
	groupArticles := g.Group(path)
	fe := &Handler{db, new2.CreateHandler(groupArticles, db, authMid, authDescribeMid)}
	g.GET("/articles", fe.Index, authDescribeMid, authMid)

	return fe
}

func (h *Handler) Index(c echo.Context) error {
	ctx := c.Request().Context()
	p, _ := strconv.Atoi(c.QueryParam(constants.PAGE))
	ps, _ := strconv.Atoi(c.QueryParam(constants.PAGE_SIZE))
	source := c.QueryParam(constants.TARGET_SOURCE)

	hxRequest, _ := strconv.ParseBool(c.Request().Header.Get("Hx-Request"))

	partial := source == constants.TARGET_SOURCE_PARTIAL && hxRequest

	var count int64
	h.DB.WithContext(ctx).
		Model(&models.Post{}).
		Count(&count)

	var posts []models.Post
	if err := h.DB.WithContext(ctx).
		Scopes(
			db.Paginate(p, ps),
			db.OrderBy("created_at", db.Descending),
			post.Shallow(),
		).
		Find(&posts).Error; err != nil {
		return errs.Wrap(err)
	}

	if len(posts) <= 0 && p > 1 {
		appendRoute := ""
		if source == constants.TARGET_SOURCE_PARTIAL {
			appendRoute = "&source=source-partial"
		}
		path := fmt.Sprintf("/dashboard/articles?page=%d&pageSize=%d%s", p-1, ps, appendRoute)
		fmt.Println(path)
		return c.Redirect(http.StatusFound, path)
	}

	opts := pub_variables.DashboardOpts{
		Nav:         sidebar.Nav(0),
		BreadCrumbs: breadcrumb.Breadcrumbs("dashboard/articles"),
	}

	pageSizes := size.Dropdown(p, ps)
	pages := page.Buttons(tern.Int(p, 1), tern.Int(ps, 10), int(count))
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
