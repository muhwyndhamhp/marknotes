package admin

import (
	"fmt"
	"github.com/muhwyndhamhp/marknotes/internal"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/admin/contact"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/admin/index"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/admin/unauthorized"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/common/variables"
	"github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/jwt"
)

type handler struct {
	app *internal.Application
}

func NewHandler(
	g *echo.Group,
	app *internal.Application,
) {
	fe := &handler{app: app}

	g.GET("", fe.Index, app.DescribeAuthWare, app.CacheControlWare)
	g.GET("/unauthorized", fe.Unauthorized, app.CacheControlWare)
	g.GET("/resume", fe.Resume, app.CacheControlWare)
	g.GET("/contact", fe.Contact, app.DescribeAuthWare, app.CacheControlWare)
}

func (fe *handler) Contact(c echo.Context) error {
	userID := jwt.AppendAndReturnUserID(c, map[string]interface{}{})

	bodyOpts := variables.BodyOpts{
		HeaderButtons: AppendHeaderButtons(userID),
		Component:     nil,
	}

	return template.AssertRender(c, http.StatusOK, contact.Contact(bodyOpts))
}

func (fe *handler) Resume(c echo.Context) error {
	return c.Redirect(http.StatusFound,
		fmt.Sprintf("/posts/%s", config.Get(config.RESUME_POST_ID)))
}

func (fe *handler) Index(c echo.Context) error {
	userID := jwt.AppendAndReturnUserID(c, map[string]interface{}{})

	bodyOpts := variables.BodyOpts{
		HeaderButtons: AppendHeaderButtons(userID),
		Component:     nil,
		HideTitle:     true,
	}

	index := index.Index(bodyOpts)

	return template.AssertRender(c, http.StatusOK, index)
}

func (fe *handler) Unauthorized(c echo.Context) error {
	bodyOpts := variables.BodyOpts{
		HeaderButtons: AppendHeaderButtons(0),
		Component:     nil,
	}

	return template.AssertRender(
		c,
		http.StatusOK,
		unauthorized.Unauthorized(bodyOpts),
	)
}
