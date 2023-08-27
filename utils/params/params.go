package params

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/utils/constants"
)

func GetCommonParams(c echo.Context) (
	page int,
	pageSize int,
	sortBy string,
	status string,
) {
	page, _ = strconv.Atoi(c.QueryParam(constants.PAGE))
	pageSize, _ = strconv.Atoi(c.QueryParam(constants.PAGE_SIZE))
	sortBy = c.QueryParam(constants.SORT_BY)
	status = c.QueryParam(constants.STATUS)

	return page, pageSize, sortBy, status
}
