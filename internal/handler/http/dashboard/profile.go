package dashboard

import (
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/dashboard/profile"
	"net/http"

	"github.com/labstack/echo/v4"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	templates "github.com/muhwyndhamhp/marknotes/template"
)

func (fe *handler) Profile(c echo.Context) error {
	opts := pub_variables.DashboardOpts{Nav: nav(1)}

	dashboard := profile.Profile(opts)

	return templates.AssertRender(c, http.StatusOK, dashboard)
}
