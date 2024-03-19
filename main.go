package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/db"
	"github.com/muhwyndhamhp/marknotes/db/migration"
	"github.com/muhwyndhamhp/marknotes/middlewares"
	"github.com/muhwyndhamhp/marknotes/pkg/admin"
	"github.com/muhwyndhamhp/marknotes/pkg/auth"
	_userRepo "github.com/muhwyndhamhp/marknotes/pkg/auth/repository"
	"github.com/muhwyndhamhp/marknotes/pkg/dashboard"
	"github.com/muhwyndhamhp/marknotes/pkg/post"
	_postRepo "github.com/muhwyndhamhp/marknotes/pkg/post/repository"
	"github.com/muhwyndhamhp/marknotes/pkg/site"
	"github.com/muhwyndhamhp/marknotes/pkg/tag"
	_tagRepo "github.com/muhwyndhamhp/marknotes/pkg/tag/repository"
	"github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/clerkauth"
	"github.com/muhwyndhamhp/marknotes/utils/cloudbucket"
	"github.com/muhwyndhamhp/marknotes/utils/jwt"
	"github.com/muhwyndhamhp/marknotes/utils/renderfile"
	"github.com/muhwyndhamhp/marknotes/utils/routing"
	"github.com/muhwyndhamhp/marknotes/utils/rss"
)

// nolint: typecheck
func main() {
	e := echo.New()

	clerkClient := clerkauth.NewClient(config.Get(config.CLERK_SECRET))

	routing.SetupRouter(e, clerkClient.Clerk)

	e.Use(redirectHTML())
	e.Use(middlewares.SetCachePolicy())
	e.Static("/dist", "dist")
	e.Static("/assets", "public/assets")
	e.Static("/articles", "public/articles")

	e.Static("/public/sitemap", "public/sitemap")
	e.File("/robots.txt", "public/assets/robots.txt")
	e.File("/sitemap.xml", "public/sitemap/sitemap.xml")

	template.NewTemplateRenderer(e)

	adminGroup := e.Group("")

	bucket := cloudbucket.NewS3Client()
	postRepo := _postRepo.NewPostRepository(db.GetLibSQLDB())
	userRepo := _userRepo.NewUserRepository(db.GetLibSQLDB())
	tagRepo := _tagRepo.NewTagRepository(db.GetLibSQLDB())
	htmxMid := middlewares.HTMXRequest()

	ctx := context.Background()
	err := rss.GenerateRSS(ctx, postRepo)
	if err != nil {
		panic(err)
	}
	e.File("/rss.xml", "public/assets/rss.xml")

	service := jwt.Service{SecretKey: []byte(config.Get(config.JWT_SECRET))}
	authMid := clerkClient.AuthMiddleware(userRepo)
	authDescMid := service.AuthDescribeMiddleware()
	byIDMid := middlewares.ByIDMiddleware()

	admin.NewAdminFrontend(adminGroup, postRepo, authDescMid)
	post.NewPostFrontend(adminGroup, postRepo, bucket, htmxMid, authMid, authDescMid, byIDMid)
	dashboard.NewDashboardFrontend(adminGroup, postRepo, tagRepo, htmxMid, authMid, authDescMid, byIDMid)
	tag.NewTagFrontend(adminGroup, tagRepo, authMid)
	auth.NewAuthService(adminGroup, service, config.Get(config.OAUTH_AUTHORIZE_URL),
		config.Get(config.OAUTH_ACCESSTOKEN_URL),
		config.Get(config.OAUTH_CLIENTID),
		config.Get(config.OAUTH_SECRET),
		config.Get(config.OAUTH_URL),
		userRepo)

	go func() {
		if config.Get(config.ENV) == "dev" {
			return
		}

		migration.Migrate(db.GetLibSQLDB())
	}()

	go func() {
		if config.Get(config.ENV) == "dev" {
			return
		}

		ctx := context.Background()
		_, err := bucket.UploadStatic(ctx, "dist/tailwind.css", "text/css")
		if err != nil {
			e.Logger.Fatal(err)
		}

		_, err = bucket.UploadStatic(ctx, "dist/main.js", "text/javascript")
		if err != nil {
			e.Logger.Fatal(err)
		}

		_, err = bucket.UploadStatic(ctx, "dist/editor.js", "text/javascript")
		if err != nil {
			e.Logger.Fatal(err)
		}
	}()

	go func() {
		renderfile.RenderPosts(context.Background(), postRepo)
		if config.Get(config.ENV) != "dev" {
			site.PingSitemap(postRepo)
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
