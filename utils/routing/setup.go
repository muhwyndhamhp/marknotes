package routing

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/apsystole/log"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/muhwyndhamhp/marknotes/utils/resp"
	"github.com/muhwyndhamhp/marknotes/utils/validate"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

func SetupRouter(e *echo.Echo) {
	e.HTTPErrorHandler = httpErrorHandler

	e.Pre(middleware.RemoveTrailingSlash())
	// e.Use(middleware.Recover())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log.Printf("[PANIC] %v", err)

			for _, line := range bytes.Split(stack, []byte("\n")) {
				cleaned := strings.ReplaceAll(string(line), "\t", "    ") // replace tab with 4 spaces
				cleaned = strings.TrimSpace(cleaned)                      // optional: remove leading/trailing whitespace
				log.Printf(" >> %s", cleaned)
			}

			return nil
		},
	}))

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*://localhost:*", "*://www.github.com", "*://github.com", "*.fly.dev", "*://mwyndham.dev", "unpkg.com", "cdn.jsdelivr.net", "static.cloudflare.com", "static.cloudflareinsights.com", "github.com", "*.mwyndham.dev"},
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
		_ = c.JSON(code, errors.WithStack(err))
	} else {
		log.Error(err)
		_ = resp.HTTPServerError(c)
	}
}
