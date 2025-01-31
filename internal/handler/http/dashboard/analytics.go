package dashboard

import (
	"fmt"
	"github.com/muhwyndhamhp/marknotes/internal"
	"net/http"

	"github.com/labstack/echo/v4"
	analytic "github.com/muhwyndhamhp/marknotes/analytics"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/dashboard/analytics"
	"github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"github.com/muhwyndhamhp/marknotes/utils/typeext"
)

func (fe *handler) Analytics(c echo.Context) error {
	slug := c.Param("slug")

	if slug == "" {
		return c.HTML(http.StatusOK, "")
	}

	var data []typeext.JSONB
	if err := fe.App.DB.
		WithContext(c.Request().Context()).
		Model(&internal.Analytics{}).
		Scopes(internal.GetLatest(slug)).
		Pluck("data", &data).
		Error; err != nil {
		fmt.Println(errs.Wrap(err))
		return c.HTML(http.StatusOK, "")
	}

	p, err := typeext.ConvertJSONBToStruct[analytic.AnalyticsResponse](data[0])
	if err != nil {
		fmt.Println(errs.Wrap(err))
		return c.HTML(http.StatusOK, "")
	}

	resp := analytics.Analytics(&p)

	return template.AssertRender(c, http.StatusOK, resp)
}
