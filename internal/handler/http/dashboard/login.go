package dashboard

import (
	pub_variables "github.com/muhwyndhamhp/marknotes/internal/handler/http/common/variables"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/dashboard/login"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/template"
)

func (fe *handler) Login(c echo.Context) error {
	opts := pub_variables.DashboardOpts{
		Nav:         nav(0),
		BreadCrumbs: fe.Breadcrumbs("dashboard/articles"),
	}

	loginVM := login.LoginViewModel{
		Opts: opts,
	}

	l := login.Login(&loginVM)

	return template.AssertRender(c, http.StatusOK, l)
}
