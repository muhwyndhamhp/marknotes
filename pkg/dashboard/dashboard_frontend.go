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

	g.GET("/dashboard/articles", fe.Articles, authMid)
	g.POST("/dashboard/articles/push", fe.ArticlesPush, authMid)
	g.GET("/dashboard/articles/new", fe.ArticlesNew, authMid)
	g.GET("/dashboard/articles/:id/edit", fe.ArticlesEdit, authMid)
	g.GET("/dashboard/profile", fe.Profile, authMid)
	g.GET("/dashboard/editor", fe.Editor, authMid)
	g.GET("/dashboard/tags", fe.Tags, authMid)
}

type ArticlesCreateRequest struct {
	Title   string `json:"title" validate:"required" form:"title"`
	Content string `json:"content" validate:"required" form:"content"`
	Tags    string `json:"tags" form:"tags"`
}
