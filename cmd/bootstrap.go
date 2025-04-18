package cmd

import (
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/analytics"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/db"
	"github.com/muhwyndhamhp/marknotes/internal"
	internalClerk "github.com/muhwyndhamhp/marknotes/internal/clerk"
	"github.com/muhwyndhamhp/marknotes/internal/middlewares"
	_postRepo "github.com/muhwyndhamhp/marknotes/internal/post"
	"github.com/muhwyndhamhp/marknotes/internal/renderfile"
	_tagRepo "github.com/muhwyndhamhp/marknotes/internal/tag"
	_userRepo "github.com/muhwyndhamhp/marknotes/internal/user"
	"github.com/muhwyndhamhp/marknotes/utils/cloudbucket"
	"github.com/muhwyndhamhp/marknotes/utils/imageprocessing"
	"time"
)

func Bootstrap() *internal.Application {
	app := &internal.Application{
		UserRepository:      _userRepo.NewUserRepository(db.GetLibSQLDB()),
		PostRepository:      _postRepo.NewPostRepository(db.GetLibSQLDB()),
		TagRepository:       _tagRepo.NewTagRepository(db.GetLibSQLDB()),
		AnalyticsRepository: nil,
		DB:                  db.GetLibSQLDB(),
		Bucket:              cloudbucket.NewS3Client(imageprocessing.NewClient()),
		AnalyticsClient: analytics.NewClient(
			config.Get(config.CF_ACCOUNT_ID),
			config.Get(config.CF_SERVICE_ID),
			config.Get(config.CF_ANALYTICS_GQL_API_KEY),
			config.Get(config.CF_ANALYTICS_EMAIL),
		),
	}

	app.RenderClient = renderfile.NewRenderClient(app)

	app.ClerkClient = internalClerk.NewClient(config.Get(config.CLERK_SECRET), app)
	app.RequireAuthWare = app.ClerkClient.AuthMiddleware()
	app.DescribeAuthWare = echo.WrapMiddleware(clerk.WithSessionV2(app.ClerkClient.GetClerk(), clerk.WithLeeway(60*time.Second)))
	app.CacheControlWare = middlewares.SetCachePolicy()
	app.GetIdParamWare = middlewares.ByIDMiddleware()
	app.FromHTMXRequestWare = middlewares.HTMXRequest()

	return app
}
