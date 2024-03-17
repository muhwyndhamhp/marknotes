package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/labstack/echo/v4"
)

func RegisterClient(secretKey string) {
	clerk.SetKey(secretKey)
}

func WithHeaderAuth() echo.MiddlewareFunc {
	return echo.WrapMiddleware(clerkhttp.WithHeaderAuthorization(clerkhttp.Leeway(10 * time.Second)))
}

func SessionRedirectMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ses, ok := clerk.SessionClaimsFromContext(c.Request().Context())
			if !ok {
				fmt.Println("Redirecting to login")
				err := c.Redirect(http.StatusFound, "/dashboard/login?redirect="+c.Request().URL.String())
				if err != nil {
					return err
				}
			}

			js, _ := json.MarshalIndent(ses, "", "  ")
			fmt.Println(string(js))

			return next(c)
		}
	}
}

func GetSession(c echo.Context) (*clerk.SessionClaims, error) {
	ses, ok := clerk.SessionClaimsFromContext(c.Request().Context())
	if !ok {
		return nil, c.Redirect(http.StatusFound, "/dashboard/login?redirect="+c.Request().URL.String())
	}

	return ses, nil
}
