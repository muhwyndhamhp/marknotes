package dashboard

import (
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
)

type DashboardFrontend struct {
	PostRepo models.PostRepository
	TagRepo  models.TagRepository
}

func NewDashboardFrontend(
	g *echo.Echo,
	PostRepo models.PostRepository,
	TagRepo models.TagRepository,
	authMid echo.MiddlewareFunc,
	sessionHandlerMid echo.MiddlewareFunc,
) {
	fe := &DashboardFrontend{PostRepo, TagRepo}

	g.GET("/dashboard", func(c echo.Context) error {
		return c.Redirect(301, "/dashboard/articles")
	}, sessionHandlerMid,
	// authMid
	)
	g.GET("/dashboard/articles", fe.Articles, sessionHandlerMid) // authMid

	g.POST("/dashboard/articles/push", fe.ArticlesPush, sessionHandlerMid) // authMid

	g.GET("/dashboard/articles/new", fe.ArticlesNew, sessionHandlerMid) // authMid

	g.GET("/dashboard/articles/:id", func(c echo.Context) error {
		return c.Redirect(301, "/dashboard/articles/"+c.Param("id")+"/edit")
	}, sessionHandlerMid,
	// authMid
	)
	g.GET("/dashboard/articles/:id/edit", fe.ArticlesEdit, sessionHandlerMid) // authMid

	g.GET("/dashboard/profile", fe.Profile, sessionHandlerMid) // authMid

	g.GET("/dashboard/editor", fe.Editor, sessionHandlerMid) // authMid

	g.GET("/dashboard/tags", fe.Tags, sessionHandlerMid) // authMid

	g.GET("/dismiss", func(c echo.Context) error {
		// return empty html
		return c.HTML(200, "")
	})
	g.GET("/dashboard/load-iframe", fe.LoadIframe, sessionHandlerMid) // authMid

	g.GET("/dashboard/login", fe.Login)
}

type ArticlesCreateRequest struct {
	Title       string `json:"title" validate:"required" form:"title"`
	Content     string `json:"content" validate:"required" form:"content"`
	Tags        string `json:"tags" form:"tags"`
	ContentJSON string `json:"content_json" form:"content_json" validate:"required"`
}
