package dashboard

import (
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/internal"
)

type handler struct {
	App *internal.Application
}

func NewHandler(
	g *echo.Group,
	app *internal.Application,
) {
	fe := &handler{app}

	g.GET(
		"/dashboard",
		func(c echo.Context) error {
			return c.Redirect(301, "/dashboard/articles")
		},
		app.RequireAuthWare,
	)

	g.GET("/dashboard/articles", fe.Articles, app.RequireAuthWare)
	g.POST("/dashboard/articles/push", fe.ArticlesPush, app.RequireAuthWare)
	g.GET("/dashboard/articles/new", fe.ArticlesNew, app.RequireAuthWare, app.CacheControlWare)

	g.GET(
		"/dashboard/articles/:id",
		func(c echo.Context) error { return c.Redirect(301, "/dashboard/articles/"+c.Param("id")+"/edit") },

		app.RequireAuthWare,
	)

	g.GET("/dismiss", func(c echo.Context) error { return c.HTML(200, "") })

	g.GET("/dashboard/articles/:id/edit", fe.ArticlesEdit, app.RequireAuthWare)
	g.PUT("/dashboard/articles/:id/delete", fe.ArticlesDelete, app.RequireAuthWare)
	g.GET("/dashboard/profile", fe.Profile, app.RequireAuthWare)
	g.GET("/dashboard/editor", fe.Editor, app.RequireAuthWare)
	g.GET("/dashboard/tags", fe.Tags, app.RequireAuthWare)
	g.GET("/dashboard/articles/:id/export/html", fe.ExportHTML, app.RequireAuthWare)
	g.GET("/dashboard/articles/:id/export/json", fe.ExportJSON, app.RequireAuthWare)
	g.GET("/dashboard/articles/:id/export/markdown", fe.ExportMarkdown, app.RequireAuthWare)
	g.GET("/dashboard/load-iframe", fe.LoadIframe)
	g.GET("/dashboard/login", fe.Login)
	g.GET("/dashboard/analytics/:slug", fe.Analytics, app.RequireAuthWare)
}

type ArticlesCreateRequest struct {
	Title           string `json:"title" validate:"required" form:"title"`
	Content         string `json:"content" validate:"required" form:"content"`
	Tags            string `json:"tags" form:"tags"`
	ContentJSON     string `json:"content_json" form:"content_json" validate:"required"`
	HeaderImageURL  string `json:"header_image_url" form:"header_image_url"`
	MarkdownContent string `json:"markdown_content" form:"markdown_content"`
}
