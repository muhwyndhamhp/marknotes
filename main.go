package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/db"
	"github.com/muhwyndhamhp/marknotes/middlewares"
	"github.com/muhwyndhamhp/marknotes/pkg/admin"
	"github.com/muhwyndhamhp/marknotes/pkg/auth"
	_userRepo "github.com/muhwyndhamhp/marknotes/pkg/auth/repository"
	"github.com/muhwyndhamhp/marknotes/pkg/dashboard"
	"github.com/muhwyndhamhp/marknotes/pkg/post"
	_postRepo "github.com/muhwyndhamhp/marknotes/pkg/post/repository"
	"github.com/muhwyndhamhp/marknotes/pkg/tag"
	_tagRepo "github.com/muhwyndhamhp/marknotes/pkg/tag/repository"
	"github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/jwt"
	"github.com/muhwyndhamhp/marknotes/utils/routing"
)

// nolint: typecheck
func main() {
	e := echo.New()
	routing.SetupRouter(e)

	rg := e.Group("")
	rg.Use(middlewares.SetCachePolicy())
	rg.Static("/dist", "dist")
	rg.Static("/assets", "public/assets")

	e.Static("/public/sitemap", "public/sitemap")
	e.File("/robots.txt", "public/assets/robots.txt")
	e.File("/sitemap.xml", "public/sitemap/sitemap.xml")

	template.NewTemplateRenderer(e)

	adminGroup := e.Group("")

	postRepo := _postRepo.NewPostRepository(db.GetDB())
	userRepo := _userRepo.NewUserRepository(db.GetDB())
	tagRepo := _tagRepo.NewTagRepository(db.GetDB())
	htmxMid := middlewares.HTMXRequest()

	service := jwt.Service{SecretKey: []byte(config.Get(config.JWT_SECRET))}
	authMid := service.AuthMiddleware()
	authDescMid := service.AuthDescribeMiddleware()
	byIDMid := middlewares.ByIDMiddleware()

	admin.NewAdminFrontend(adminGroup, postRepo, authDescMid)
	post.NewPostFrontend(adminGroup, postRepo, htmxMid, authMid, authDescMid, byIDMid)
	dashboard.NewDashboardFrontend(adminGroup, postRepo, htmxMid, authMid, authDescMid, byIDMid)
	tag.NewTagFrontend(adminGroup, tagRepo, authMid)
	auth.NewAuthService(adminGroup, service, config.Get(config.OAUTH_AUTHORIZE_URL),
		config.Get(config.OAUTH_ACCESSTOKEN_URL),
		config.Get(config.OAUTH_CLIENTID),
		config.Get(config.OAUTH_SECRET),
		config.Get(config.OAUTH_URL),
		userRepo)

	// go func() {
	// 	migration.Migrate(db.GetDB())
	// }()
	//
	// go func() {
	// 	site.PingSitemap(postRepo)
	// }()

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Get(config.APP_PORT))))
}
