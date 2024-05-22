package middlewares

import "github.com/labstack/echo/v4"

func SetCachePolicy() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "max-age=172800")
			return next(c)
		}
	}
}
