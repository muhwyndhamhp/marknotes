package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/gotes-mx/config"
	"github.com/muhwyndhamhp/gotes-mx/db"
	"github.com/muhwyndhamhp/gotes-mx/pkg/admin"
	"github.com/muhwyndhamhp/gotes-mx/pkg/repository"
	"github.com/muhwyndhamhp/gotes-mx/template"
	"github.com/muhwyndhamhp/gotes-mx/utils/routing"
)

func main() {

	e := echo.New()
	routing.SetupRouter(e)

	e.Static("/dist", "dist")

	template.NewTemplateRenderer(e,
		"public/*.html",
		"public/admin/views/*.html",
		"public/admin/components/*.html",
	)

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", "This is example how templating works!")
	})

	adminGroup := e.Group("admin")

	postRepo := repository.NewPostRepository(db.GetDB())
	admin.NewAdminFrontend(adminGroup, postRepo)
	admin.NewPostFrontend(adminGroup, postRepo)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Get(config.APP_PORT))))
}
