package internal

import (
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/analytics"
	"github.com/muhwyndhamhp/marknotes/utils/clerkauth"
	"github.com/muhwyndhamhp/marknotes/utils/cloudbucket"
	"github.com/muhwyndhamhp/marknotes/utils/renderfile"
	"gorm.io/gorm"
)

type Application struct {
	// Repositories
	UserRepository      UserRepository
	PostRepository      PostRepository
	TagRepository       TagRepository
	AnalyticsRepository AnalyticsRepository

	// Internal Plumbings and Clients
	DB              *gorm.DB
	ClerkClient     *clerkauth.Client
	Bucket          *cloudbucket.S3Client
	RenderClient    *renderfile.RenderClient
	AnalyticsClient *analytics.Client

	// Middlewares
	RequireAuthWare     echo.MiddlewareFunc
	DescribeAuthWare    echo.MiddlewareFunc
	CacheControlWare    echo.MiddlewareFunc
	GetIdParamWare      echo.MiddlewareFunc
	FromHTMXRequestWare echo.MiddlewareFunc
}
