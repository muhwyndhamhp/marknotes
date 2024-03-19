package routing

import (
	"net/http"
	"time"

	"github.com/apsystole/log"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/muhwyndhamhp/marknotes/utils/resp"
	"github.com/muhwyndhamhp/marknotes/utils/validate"
	"golang.org/x/time/rate"
)

func SetupRouter(e *echo.Echo, clerkClient clerk.Client) {
	e.HTTPErrorHandler = httpErrorHandler

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(20),
	)))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*://localhost:*", "*://www.github.com", "*://github.com", "*.fly.dev", "*://mwyndham.dev", "unpkg.com", "cdn.jsdelivr.net", "static.cloudflare.com", "static.cloudflareinsights.com", "github.com"},
	}))

	sessionMid := echo.WrapMiddleware(clerk.WithSessionV2(clerkClient, clerk.WithLeeway(10*time.Second)))
	e.Use(sessionMid)

	e.Validator = &validate.CustomValidator{
		Validator: validator.New(),
	}
}

func httpErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	if code != http.StatusInternalServerError {
		_ = c.JSON(code, err)
	} else {
		log.Error(err)
		_ = resp.HTTPServerError(c)
	}
}
