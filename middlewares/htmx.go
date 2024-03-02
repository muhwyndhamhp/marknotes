package middlewares

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func HTMXRequest() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			hx_request, err := strconv.ParseBool(c.Request().Header.Get("Hx-Request"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, nil)
			}

			if !hx_request {
				return c.JSON(http.StatusBadRequest, nil)
			}
			return next(c)
		}
	}
}
