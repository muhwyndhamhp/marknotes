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
	ReplyRepository     ReplyRepository

	// Internal Plumbings and Clients
	DB              *gorm.DB
	Bucket          *cloudbucket.S3Client
	RenderClient    RenderFile
	AnalyticsClient *analytics.Client
	OpenAuth        OpenAuth
	ProfanityCheck  ProfanityCheck
	LLM             LLM

	// Middlewares
	RequireAuthWare     echo.MiddlewareFunc
	CacheControlWare    echo.MiddlewareFunc
	GetIdParamWare      echo.MiddlewareFunc
	FromHTMXRequestWare echo.MiddlewareFunc
}
