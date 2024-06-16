package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/analytics"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/db"
	"github.com/muhwyndhamhp/marknotes/middlewares"
	"github.com/muhwyndhamhp/marknotes/pkg/admin"
	"github.com/muhwyndhamhp/marknotes/pkg/auth"
	_userRepo "github.com/muhwyndhamhp/marknotes/pkg/auth/repository"
	"github.com/muhwyndhamhp/marknotes/pkg/dashboard"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/pkg/post"
	_postRepo "github.com/muhwyndhamhp/marknotes/pkg/post/repository"
	"github.com/muhwyndhamhp/marknotes/pkg/site"
	_tagRepo "github.com/muhwyndhamhp/marknotes/pkg/tag/repository"
	"github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/clerkauth"
	"github.com/muhwyndhamhp/marknotes/utils/cloudbucket"
	"github.com/muhwyndhamhp/marknotes/utils/imageprocessing"
	"github.com/muhwyndhamhp/marknotes/utils/jwt"
	"github.com/muhwyndhamhp/marknotes/utils/renderfile"
	"github.com/muhwyndhamhp/marknotes/utils/routing"
	"github.com/muhwyndhamhp/marknotes/utils/rss"
	_webDashboard "github.com/muhwyndhamhp/marknotes/web/routes/dashboard"
)

// nolint: typecheck
func main() {
	e := echo.New()

	clerkClient := clerkauth.NewClient(config.Get(config.CLERK_SECRET))

	routing.SetupRouter(e, clerkClient.Clerk)

	e.Use(redirectHTML())
	// e.Use(middlewares.SetCachePolicy())
	e.Static("/dist", "dist")
	e.Static("/assets", "public/assets")
	e.Static("/articles", config.Get(config.POST_RENDER_PATH))

	e.Static("/public/sitemap", "public/sitemap")
	e.File("/robots.txt", "public/assets/robots.txt")
	e.File("/sitemap.xml", "public/sitemap/sitemap.xml")

	template.NewTemplateRenderer(e)

	adminGroup := e.Group("")

	iproc := imageprocessing.NewClient()
	bucket := cloudbucket.NewS3Client(iproc)
	postRepo := _postRepo.NewPostRepository(db.GetLibSQLDB())
	userRepo := _userRepo.NewUserRepository(db.GetLibSQLDB())
	tagRepo := _tagRepo.NewTagRepository(db.GetLibSQLDB())
	htmxMid := middlewares.HTMXRequest()

	analyticsClient := analytics.NewClient(
		config.Get(config.CF_ACCOUNT_ID),
		config.Get(config.CF_SERVICE_ID),
		config.Get(config.CF_ANALYTICS_GQL_API_KEY),
		config.Get(config.CF_ANALYTICS_EMAIL),
	)

	ctx := context.Background()
	err := rss.GenerateRSS(ctx, postRepo)
	if err != nil {
		panic(err)
	}
	e.File("/rss.xml", "public/assets/rss.xml")

	service := jwt.Service{SecretKey: []byte(config.Get(config.JWT_SECRET))}
	authMid := clerkClient.AuthMiddleware(userRepo)
	authDescMid := echo.WrapMiddleware(clerk.WithSessionV2(clerkClient.Clerk, clerk.WithLeeway(60*time.Second)))
	byIDMid := middlewares.ByIDMiddleware()
	cacheControlMid := middlewares.SetCachePolicy()

	e.GET("/touch", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}, authDescMid, authMid)

	admin.NewAdminFrontend(adminGroup, postRepo, authDescMid, cacheControlMid)
	post.NewPostFrontend(adminGroup, postRepo, bucket, htmxMid, authMid, authDescMid, byIDMid, cacheControlMid)
	dashboard.NewDashboardFrontend(
		adminGroup,
		db.GetLibSQLDB(),
		postRepo, userRepo, tagRepo,
		clerkClient,
		htmxMid, authMid, authDescMid, byIDMid,
		bucket,
		cacheControlMid,
	)

	_webDashboard.NewHandler(adminGroup, db.GetLibSQLDB(), authMid, authDescMid)

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

		ctx := context.Background()
		_, err := bucket.UploadStatic(ctx, "dist/tailwind_v4.css", "", "text/css")
		if err != nil {
			e.Logger.Fatal(err)
		}

		_, err = bucket.UploadStatic(ctx, "dist/main.js", "", "text/javascript")
		if err != nil {
			e.Logger.Fatal(err)
		}

		_, err = bucket.UploadStatic(ctx, "dist/htmx.js", "", "text/javascript")
		if err != nil {
			e.Logger.Fatal(err)
		}

		_, err = bucket.UploadStatic(ctx, "dist/auth.js", "", "text/javascript")
		if err != nil {
			e.Logger.Fatal(err)
		}

		_, err = bucket.UploadStatic(ctx, "dist/editor.js", "", "text/javascript")
		if err != nil {
			e.Logger.Fatal(err)
		}
	}()

	go func() {
		ctx := context.Background()
		renderfile.RenderPosts(ctx, postRepo, bucket)
		if config.Get(config.ENV) != "dev" {
			site.PingSitemap(postRepo)
		}

		renderfile.RenderMarkdowns(ctx, postRepo)
	}()

	go func() {
		for {
			ctx, cancel := context.WithCancel(context.Background())
			if err := models.CacheAnalytics(ctx, db.GetLibSQLDB(), analyticsClient); err != nil {
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
