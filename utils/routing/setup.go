package routing

import (
	"net/http"

	"github.com/apsystole/log"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/muhwyndhamhp/marknotes/utils/resp"
	"github.com/muhwyndhamhp/marknotes/utils/validate"
	"golang.org/x/time/rate"
)

func SetupRouter(e *echo.Echo) {
	e.HTTPErrorHandler = httpErrorHandler

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(20),
	)))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*://localhost:*", "*://www.github.com", "*://github.com", "*.fly.dev", "*://mwyndham.dev"},
	}))

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
