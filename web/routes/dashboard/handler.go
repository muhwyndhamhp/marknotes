package dashboard

import (
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/web/routes/dashboard/articles"
	"gorm.io/gorm"
	"net/http"
)

type Handler struct {
	Articles *articles.Handler
}

func NewHandler(
	g *echo.Group,
	db *gorm.DB,
	authMid echo.MiddlewareFunc,
	authDescribeMid echo.MiddlewareFunc,
) *Handler {
	path := "/dashboard"

	dashGroup := g.Group(path)

	fe := &Handler{Articles: articles.NewHandler(dashGroup, db, authMid, authDescribeMid)}

	g.GET(path, fe.ToArticles, authDescribeMid, authMid)

	return fe
}

func (h *Handler) ToArticles(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/dashboard/articles")
}
