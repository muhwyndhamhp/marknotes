package cmd

import (
	"github.com/muhwyndhamhp/marknotes/analytics"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/db"
	"github.com/muhwyndhamhp/marknotes/internal"
	"github.com/muhwyndhamhp/marknotes/internal/middlewares"
	"github.com/muhwyndhamhp/marknotes/internal/openauth"
	_postRepo "github.com/muhwyndhamhp/marknotes/internal/post"
	"github.com/muhwyndhamhp/marknotes/internal/profanity"
	"github.com/muhwyndhamhp/marknotes/internal/renderfile"
	_replyRepo "github.com/muhwyndhamhp/marknotes/internal/reply"
	_tagRepo "github.com/muhwyndhamhp/marknotes/internal/tag"
	_userRepo "github.com/muhwyndhamhp/marknotes/internal/user"
	"github.com/muhwyndhamhp/marknotes/utils/cloudbucket"
	"github.com/muhwyndhamhp/marknotes/utils/imageprocessing"
)

func Bootstrap() *internal.Application {
	app := &internal.Application{
		UserRepository:      _userRepo.NewUserRepository(db.GetLibSQLDB()),
		PostRepository:      _postRepo.NewPostRepository(db.GetLibSQLDB()),
		TagRepository:       _tagRepo.NewTagRepository(db.GetLibSQLDB()),
		ReplyRepository:     _replyRepo.NewRepository(db.GetLibSQLDB()),
		AnalyticsRepository: nil,
		DB:                  db.GetLibSQLDB(),
		Bucket:              cloudbucket.NewS3Client(imageprocessing.NewClient()),
		ProfanityCheck:      profanity.NewClient(),
		AnalyticsClient: analytics.NewClient(
			config.Get(config.CF_ACCOUNT_ID),
			config.Get(config.CF_SERVICE_ID),
			config.Get(config.CF_ANALYTICS_GQL_API_KEY),
			config.Get(config.CF_ANALYTICS_EMAIL),
		),
	}

	app.RenderClient = renderfile.NewRenderClient(app)

	app.CacheControlWare = middlewares.SetCachePolicy()
	app.GetIdParamWare = middlewares.ByIDMiddleware()
	app.FromHTMXRequestWare = middlewares.HTMXRequest()

	app.OpenAuth = openauth.NewClient(app)
	app.RequireAuthWare = app.OpenAuth.AuthMiddleware()

	return app
}
