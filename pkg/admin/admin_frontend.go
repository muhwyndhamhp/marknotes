package admin

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	"github.com/muhwyndhamhp/marknotes/utils/jwt"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
)

type AdminFrontend struct {
	repo models.PostRepository
}

func NewAdminFrontend(
	g *echo.Group,
	repo models.PostRepository,
	authDescMid echo.MiddlewareFunc,
) {
	fe := &AdminFrontend{
		repo: repo,
	}

	g.GET("", fe.Index, authDescMid)
	g.GET("/unauthorized", fe.Unauthorized)
	g.GET("/resume", fe.Resume)
	g.GET("/contact", fe.Contact, authDescMid)
}

func (fe *AdminFrontend) Contact(c echo.Context) error {
	resp := map[string]interface{}{}

	userID := jwt.AppendAndReturnUserID(c, resp)
	resp[HeaderButtonsKey] = AppendHeaderButtons(userID)

	return c.Render(http.StatusOK, "contact", resp)
}

func (fe *AdminFrontend) Resume(c echo.Context) error {
	return c.Redirect(http.StatusFound,
		fmt.Sprintf("/posts/%s", config.Get(config.RESUME_POST_ID)))
}

func (fe *AdminFrontend) Index(c echo.Context) error {
	ctx := c.Request().Context()

	posts, err := fe.repo.Get(ctx,
		scopes.Paginate(1, 5),
		scopes.OrderBy("published_at", scopes.Descending),
		scopes.WithStatus(values.Published),
	)

	if err != nil {
		return err
	}
	resp := map[string]interface{}{
		"Posts": posts,
	}

	userID := jwt.AppendAndReturnUserID(c, resp)
	resp[HeaderButtonsKey] = AppendHeaderButtons(userID)

	return c.Render(http.StatusOK, "index", resp)
}

func (fe *AdminFrontend) Unauthorized(c echo.Context) error {
	return c.Render(http.StatusOK, "unauthorized", nil)
}
