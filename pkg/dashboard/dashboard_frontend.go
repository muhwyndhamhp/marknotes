package dashboard

import (
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/utils/clerkauth"
	"github.com/muhwyndhamhp/marknotes/utils/cloudbucket"
)

type DashboardFrontend struct {
	PostRepo    models.PostRepository
	TagRepo     models.TagRepository
	UserRepo    models.UserRepository
	ClerkClient *clerkauth.Client
	Bucket      *cloudbucket.S3Client
}

func NewDashboardFrontend(
	g *echo.Group,
	PostRepo models.PostRepository,
	UserRepo models.UserRepository,
	TagRepo models.TagRepository,
	ClerkClient *clerkauth.Client,
	htmxMid echo.MiddlewareFunc,
	authMid echo.MiddlewareFunc,
	authDescribeMid echo.MiddlewareFunc,
	byIDMiddleware echo.MiddlewareFunc,
	bucket *cloudbucket.S3Client,
	cacheControlMid echo.MiddlewareFunc,
) {
	fe := &DashboardFrontend{PostRepo, TagRepo, UserRepo, ClerkClient, bucket}

	g.GET("/dashboard", func(c echo.Context) error {
		return c.Redirect(301, "/dashboard/articles")
	}, authDescribeMid, authMid)
	g.GET("/dashboard/articles", fe.Articles, authDescribeMid, authMid)
	g.POST("/dashboard/articles/push", fe.ArticlesPush, authDescribeMid, authMid)
	g.GET("/dashboard/articles/new", fe.ArticlesNew, authDescribeMid, authMid, cacheControlMid)
	g.GET("/dashboard/articles/:id", func(c echo.Context) error {
		return c.Redirect(301, "/dashboard/articles/"+c.Param("id")+"/edit")
	}, authDescribeMid, authMid)
	g.GET("/dashboard/articles/:id/edit", fe.ArticlesEdit, authDescribeMid, authMid)
	g.PUT("/dashboard/articles/:id/delete", fe.ArticlesDelete, authDescribeMid, authMid)
	g.GET("/dashboard/profile", fe.Profile, authDescribeMid, authMid)
	g.GET("/dashboard/editor", fe.Editor, authDescribeMid, authMid)
	g.GET("/dashboard/tags", fe.Tags, authDescribeMid, authMid)
	g.GET("/dashboard/articles/:id/export/html", fe.ExportHTML, authDescribeMid, authMid)
	g.GET("/dashboard/articles/:id/export/json", fe.ExportJSON, authDescribeMid, authMid)
	g.GET("/dismiss", func(c echo.Context) error {
		// return empty html
		return c.HTML(200, "")
	})
	g.GET("/dashboard/load-iframe", fe.LoadIframe)
	g.GET("/dashboard/login", fe.Login)
}

type ArticlesCreateRequest struct {
	Title          string `json:"title" validate:"required" form:"title"`
	Content        string `json:"content" validate:"required" form:"content"`
	Tags           string `json:"tags" form:"tags"`
	ContentJSON    string `json:"content_json" form:"content_json" validate:"required"`
	HeaderImageURL string `json:"header_image_url" form:"header_image_url"`
}
