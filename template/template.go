package template

import (
	"errors"
	"fmt"
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

	writer := newWriterFromWriter(w)
	res := component.Render(c.Request().Context(), writer)
	if name == "templ-log" {
		fmt.Println(string(bytes))
	}

	bytes = []byte("")
	return res
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

func AssertRenderLog(c echo.Context, statusCode int, component templ.Component) error {
	return c.Render(statusCode, "templ-log", component)
}

type writer struct {
	existing io.Writer
}

var bytes = []byte("")

// Write implements io.Writer.
func (w writer) Write(p []byte) (n int, err error) {
	bytes = append(bytes, p...)
	return w.existing.Write(p)
}

func newWriterFromWriter(w io.Writer) io.Writer {
	return writer{w}
}
