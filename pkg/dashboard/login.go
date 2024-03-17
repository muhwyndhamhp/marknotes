package dashboard

import (
	"net/http"

	"github.com/labstack/echo/v4"
	pub_dashboard_login "github.com/muhwyndhamhp/marknotes/pub/pages/dashboards/login"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	templates "github.com/muhwyndhamhp/marknotes/template"
)

func (fe *DashboardFrontend) Login(c echo.Context) error {
	opts := pub_variables.DashboardOpts{
		Nav:         nav(0),
		BreadCrumbs: fe.Breadcrumbs("dashboard/articles"),
	}

	loginVM := pub_dashboard_login.LoginVM{
		Opts: opts,
	}

	login := pub_dashboard_login.Login(&loginVM)

	return templates.AssertRender(c, http.StatusOK, login)
}
