package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/muhwyndhamhp/marknotes/cmd"
	"github.com/muhwyndhamhp/marknotes/internal"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/admin"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/dashboard"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/openauth"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/post"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/db"
	"github.com/muhwyndhamhp/marknotes/pkg/site"
	"github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/routing"
	"github.com/muhwyndhamhp/marknotes/utils/rss"
	_ "github.com/toolbeam/openauth/client"
)

// nolint: typecheck
func main() {
	e := echo.New()

	app := cmd.Bootstrap()

	routing.SetupRouter(e)

	e.Use(redirectHTML())
	e.Static("/dist", "dist")
	e.Static("/assets", "public/assets")
	e.Static("/articles", config.Get(config.POST_RENDER_PATH))

	e.Static("/public/sitemap", "public/sitemap")
	e.File("/robots.txt", "public/assets/robots.txt")
	e.File("/sitemap.xml", "public/sitemap/sitemap.xml")

	template.NewTemplateRenderer(e)

	adminGroup := e.Group("")

	ctx := context.Background()
	err := rss.GenerateRSS(ctx, app.PostRepository)
	if err != nil {
		panic(err)
	}
	e.File("/rss.xml", "public/assets/rss.xml")

	e.GET("/touch", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}, app.RequireAuthWare)

	admin.NewHandler(adminGroup, app)
	post.NewHandler(adminGroup, app)
	dashboard.NewHandler(adminGroup, app)
	openauth.NewHandler(adminGroup, app)

	go func() {
		if config.Get(config.ENV) == "dev" {
			return
		}

		ctx := context.Background()
		_, err := app.Bucket.UploadStatic(ctx, "dist/tailwind_v4.css", "", "text/css")
		if err != nil {
			e.Logger.Fatal(err)
		}

		_, err = app.Bucket.UploadStatic(ctx, "dist/main.js", "", "text/javascript")
		if err != nil {
			e.Logger.Fatal(err)
		}

		_, err = app.Bucket.UploadStatic(ctx, "dist/htmx.js", "", "text/javascript")
		if err != nil {
			e.Logger.Fatal(err)
		}

		_, err = app.Bucket.UploadStatic(ctx, "dist/auth.js", "", "text/javascript")
		if err != nil {
			e.Logger.Fatal(err)
		}

		_, err = app.Bucket.UploadStatic(ctx, "dist/editor.js", "", "text/javascript")
		if err != nil {
			e.Logger.Fatal(err)
		}
	}()

	go func() {
		ctx := context.Background()
		app.RenderClient.RenderPosts(ctx)
		if config.Get(config.ENV) != "dev" {
			site.PingSitemap(app.PostRepository)
		}

		app.RenderClient.RenderMarkdowns(ctx)
	}()

	go func() {
		for {
			ctx, cancel := context.WithCancel(context.Background())
			if err := internal.CacheAnalytics(ctx, db.GetLibSQLDB(), app.AnalyticsClient); err != nil {
				e.Logger.Error(err)
			}

			time.Sleep(15 * time.Minute)
			cancel()
		}
	}()

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Get(config.APP_PORT))))
}

func redirectHTML() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestedPath := c.Request().URL.Path
			if strings.HasPrefix(requestedPath, "/articles/") && !strings.HasSuffix(requestedPath, ".html") {
				newPath := requestedPath + ".html"
				return c.Redirect(http.StatusMovedPermanently, newPath)
			}

			return next(c)
		}
	}
}
