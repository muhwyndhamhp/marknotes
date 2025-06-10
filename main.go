package main

import (
	"context"
	"fmt"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/replies"
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

func main() {
	ctx := context.Background()
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

	g := e.Group("")

	e.File("/rss.xml", "public/assets/rss.xml")

	admin.NewHandler(g, app)
	post.NewHandler(g, app)
	dashboard.NewHandler(g, app)
	openauth.NewHandler(g, app)
	replies.NewHandler(g, app)

	go func() {
		if err := rss.GenerateRSS(ctx, app.PostRepository); err != nil {
			e.Logger.Fatal(err)
		}
	}()

	go func() { uploadStatics(e, app) }()

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
			if config.Get(config.ENV) == "dev" {
				break
			}

			ctx, cancel := context.WithCancel(context.Background())
			if err := internal.CacheAnalytics(ctx, db.GetLibSQLDB(), app.AnalyticsClient); err != nil {
				e.Logger.Error(err)
			}

			time.Sleep(3 * time.Hour)
			cancel()
		}
	}()

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Get(config.APP_PORT))))
}

func uploadStatics(e *echo.Echo, app *internal.Application) {
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
