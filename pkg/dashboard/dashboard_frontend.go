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
	g *echo.Group,
	PostRepo models.PostRepository,
	TagRepo models.TagRepository,
	htmxMid echo.MiddlewareFunc,
	authMid echo.MiddlewareFunc,
	authDescribeMid echo.MiddlewareFunc,
	byIDMiddleware echo.MiddlewareFunc,
) {
	fe := &DashboardFrontend{PostRepo, TagRepo}

	g.GET("/dashboard", func(c echo.Context) error {
		return c.Redirect(301, "/dashboard/articles")
	}, authDescribeMid, authMid)
	g.GET("/dashboard/articles", fe.Articles, authDescribeMid, authMid)
	g.POST("/dashboard/articles/push", fe.ArticlesPush, authDescribeMid, authMid)
	g.GET("/dashboard/articles/new", fe.ArticlesNew, authDescribeMid, authMid)
	g.GET("/dashboard/articles/:id", func(c echo.Context) error {
		return c.Redirect(301, "/dashboard/articles/"+c.Param("id")+"/edit")
	}, authDescribeMid, authMid)
	g.GET("/dashboard/articles/:id/edit", fe.ArticlesEdit, authDescribeMid, authMid)
	g.GET("/dashboard/profile", fe.Profile, authDescribeMid, authMid)
	g.GET("/dashboard/editor", fe.Editor, authDescribeMid, authMid)
	g.GET("/dashboard/tags", fe.Tags, authDescribeMid, authMid)
	g.GET("/dismiss", func(c echo.Context) error {
		// return empty html
		return c.HTML(200, "")
	})
	g.GET("/dashboard/load-iframe", fe.LoadIframe, authDescribeMid, authMid)
}

type ArticlesCreateRequest struct {
	Title       string `json:"title" validate:"required" form:"title"`
	Content     string `json:"content" validate:"required" form:"content"`
	Tags        string `json:"tags" form:"tags"`
	ContentJSON string `json:"content_json" form:"content_json" validate:"required"`
}
