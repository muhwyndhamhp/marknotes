package internal

import (
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/labstack/echo/v4"
)

type ClerkClient interface {
	AuthMiddleware() echo.MiddlewareFunc
	GetUserFromSession(c echo.Context) (*User, error)
	GetUser(sc *clerk.SessionClaims) (*clerk.User, error)
	GetClerk() clerk.Client
}
