package internal

import (
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/analytics"
	"github.com/muhwyndhamhp/marknotes/utils/cloudbucket"
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
	ClerkClient     ClerkClient
	Bucket          *cloudbucket.S3Client
	RenderClient    RenderFile
	AnalyticsClient *analytics.Client

	// Middlewares
	RequireAuthWare     echo.MiddlewareFunc
	DescribeAuthWare    echo.MiddlewareFunc
	CacheControlWare    echo.MiddlewareFunc
	GetIdParamWare      echo.MiddlewareFunc
	FromHTMXRequestWare echo.MiddlewareFunc
}
