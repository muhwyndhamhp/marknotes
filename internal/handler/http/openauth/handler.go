package openauth

import (
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/internal"
)

type handler struct {
	app *internal.Application
}

func NewHandler(g *echo.Group, app *internal.Application) {
	h := &handler{
		app: app,
	}

	g.GET("/openauth/authorize", h.Authorize)
	g.GET("/openauth/callback", h.app.OpenAuth.Callback())
}

func (h *handler) Authorize(c echo.Context) error {
	url, err := h.app.OpenAuth.Authorize()
	if err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Redirect", url)

	return c.HTML(200, `<p>Redirecting to identity provider...</p>`)
}
