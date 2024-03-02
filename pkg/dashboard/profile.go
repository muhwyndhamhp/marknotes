package dashboard

import (
	"net/http"

	"github.com/labstack/echo/v4"
	pub_dashboards_profile "github.com/muhwyndhamhp/marknotes/pub/pages/dashboards/profile"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	templates "github.com/muhwyndhamhp/marknotes/template"
)

func (fe *DashboardFrontend) Profile(c echo.Context) error {
	opts := pub_variables.DashboardOpts{Nav: nav(1)}

	dashboard := pub_dashboards_profile.Profile(opts)

	return templates.AssertRender(c, http.StatusOK, dashboard)
}
