package template

import (
	"errors"
	"io"
	"text/template"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if name == "templ" {
		component, isSuccess := data.(templ.Component)
		if !isSuccess {
			return errors.New("failed to parse data as templ.Component, do you pass the correct params?")
		}

		return component.Render(c.Request().Context(), w)
	} else {
		return t.Templates.ExecuteTemplate(w, name, data)
	}
}

func NewTemplateRenderer(e *echo.Echo, paths ...string) {
	tmpl := &template.Template{}
	for i := range paths {
		template.Must(tmpl.ParseGlob(paths[i]))
	}
	t := newTemplate(tmpl)
	e.Renderer = t
}

func newTemplate(templates *template.Template) echo.Renderer {
	return &Template{
		Templates: templates,
	}
}

func AssertRender(c echo.Context, statusCode int, component templ.Component) error {
	return c.Render(statusCode, "templ", component)
}
