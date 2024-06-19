package new

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/config"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	templates "github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/web/routes/dashboard/_partials/breadcrumb"
	"github.com/muhwyndhamhp/marknotes/web/routes/dashboard/_partials/sidebar"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type Handler struct {
	DB *gorm.DB
}

func CreateHandler(
	g *echo.Group,
	db *gorm.DB,
	authMid echo.MiddlewareFunc,
	authDescribeMid echo.MiddlewareFunc,
) *Handler {
	fe := &Handler{db}

	g.GET("/new", fe.ArticlesNew, authDescribeMid, authMid)

	return fe
}

func (fe *Handler) ArticlesNew(c echo.Context) error {
	opts := pub_variables.DashboardOpts{
		Nav:         sidebar.Nav(0),
		BreadCrumbs: breadcrumb.Breadcrumbs("dashboard/articles/new"),
	}

	baseURL := strings.Split(config.Get(config.OAUTH_URL), "/callback")[0]
	uploadURL := fmt.Sprintf("%s/posts/%d/media/upload", baseURL, 0)

	vm := NewVM{
		Opts:      opts,
		UploadURL: uploadURL,
		BaseURL:   baseURL,
	}
	articlesNew := New(vm)

	return templates.AssertRender(c, http.StatusOK, articlesNew)
}
