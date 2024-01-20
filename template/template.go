package template

import (
	"errors"
	"io"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Template struct{}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	component, isSuccess := data.(templ.Component)
	if !isSuccess {
		return errors.New("failed to parse data as templ.Component, do you pass the correct params?")
	}

	return component.Render(c.Request().Context(), w)
}

func NewTemplateRenderer(e *echo.Echo) {
	t := newTemplate()
	e.Renderer = t
}

func newTemplate() echo.Renderer {
	return &Template{}
}

func AssertRender(c echo.Context, statusCode int, component templ.Component) error {
	return c.Render(statusCode, "templ", component)
}
