package resp

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/gotes-mx/utils/tern"
)

type Response struct {
	Data      interface{} `json:"data"`
	Meta      interface{} `json:"meta,omitempty"`
	Message   string      `json:"message,omitempty"`
	ErrorCode string      `json:"error_code,omitempty"`
}

func HTTPOk(c echo.Context, data interface{}) error {
	r := Response{
		Data: data,
	}

	return c.JSON(http.StatusOK, r)
}

func HTTPOkWithMeta(c echo.Context, data, meta interface{}) error {
	r := Response{
		Data: data,
		Meta: meta,
	}

	return c.JSON(http.StatusOK, r)
}

func HTTPBadRequest(c echo.Context, code, msg string) error {
	errCode := tern.String(code, ValidationError)
	r := Response{
		ErrorCode: errCode,
		Message:   msg,
	}

	return c.JSON(http.StatusBadRequest, r)
}

func HTTPForbidden(c echo.Context, code, msg string) error {
	r := Response{
		ErrorCode: code,
		Message:   msg,
	}

	return c.JSON(http.StatusForbidden, r)
}

func HTTPUnauthorized(c echo.Context) error {
	r := Response{
		ErrorCode: InvalidAuthToken,
		Message:   "Invalid or missing authorization token",
	}

	return c.JSON(http.StatusUnauthorized, r)
}

func HTTPServerError(c echo.Context) error {
	r := Response{
		ErrorCode: ServerError,
		Message:   "Something went wrong",
	}

	return c.JSON(http.StatusInternalServerError, r)
}
