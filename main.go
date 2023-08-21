package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/db"
	"github.com/muhwyndhamhp/marknotes/middlewares"
	"github.com/muhwyndhamhp/marknotes/pkg/admin"
	"github.com/muhwyndhamhp/marknotes/pkg/repository"
	"github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/routing"
)

func main() {

	e := echo.New()
	routing.SetupRouter(e)

	e.Static("/dist", "dist")
	e.Static("/assets", "public/assets")

	template.NewTemplateRenderer(e,
		"public/*.html",
		"public/admin/views/*.html",
		"public/admin/components/*.html",
		"public/styles/*.html",
	)

	// e.GET("/", func(c echo.Context) error {
	// 	return c.Render(http.StatusOK, "index", "This is example how templating works!")
	// })

	adminGroup := e.Group("")

	postRepo := repository.NewPostRepository(db.GetDB())
	htmxMid := middlewares.HTMXRequest()

	admin.NewAdminFrontend(adminGroup, postRepo)
	admin.NewPostFrontend(adminGroup, postRepo, htmxMid)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Get(config.APP_PORT))))
}
