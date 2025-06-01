package internal

import "github.com/labstack/echo/v4"

type OpenAuth interface {
	AuthMiddleware() echo.MiddlewareFunc
	GetUserFromSession(c echo.Context) (*User, error)
	Authorize() (string, error)
	Callback() func(c echo.Context) error
}
