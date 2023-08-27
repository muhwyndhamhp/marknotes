package middlewares

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

const ByIDKey = "By-ID-Key"

func ByIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id, _ := strconv.Atoi(c.Param("id"))
			if id <= 0 {
				return c.JSON(http.StatusBadRequest, nil)
			}

			c.Set(ByIDKey, id)

			return next(c)
		}
	}
}
