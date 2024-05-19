package admin

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	pub_contact "github.com/muhwyndhamhp/marknotes/pub/pages/contact"
	pub_index "github.com/muhwyndhamhp/marknotes/pub/pages/index"
	pub_unauthorized "github.com/muhwyndhamhp/marknotes/pub/pages/unauthorized"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	"github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/jwt"
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
	userID := jwt.AppendAndReturnUserID(c, map[string]interface{}{})

	bodyOpts := pub_variables.BodyOpts{
		HeaderButtons: AppendHeaderButtons(userID),
		Component:     nil,
	}

	return template.AssertRender(c, http.StatusOK, pub_contact.Contact(bodyOpts))
}

func (fe *AdminFrontend) Resume(c echo.Context) error {
	return c.Redirect(http.StatusFound,
		fmt.Sprintf("/posts/%s", config.Get(config.RESUME_POST_ID)))
}

func (fe *AdminFrontend) Index(c echo.Context) error {
	userID := jwt.AppendAndReturnUserID(c, map[string]interface{}{})

	bodyOpts := pub_variables.BodyOpts{
		HeaderButtons: AppendHeaderButtons(userID),
		Component:     nil,
		HideTitle:     true,
	}

	index := pub_index.Index(bodyOpts)

	return template.AssertRender(c, http.StatusOK, index)
}

func (fe *AdminFrontend) Unauthorized(c echo.Context) error {
	bodyOpts := pub_variables.BodyOpts{
		HeaderButtons: AppendHeaderButtons(0),
		Component:     nil,
	}

	return template.AssertRender(
		c,
		http.StatusOK,
		pub_unauthorized.Unauthorized(bodyOpts),
	)
}
