package dashboard

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/config"
)

func (fe DashboardFrontend) ExportHTML(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))

	post, err := fe.PostRepo.GetByID(ctx, uint(id))
	if err != nil {
		return err
	}

	fp := filepath.Join(config.Get(config.POST_RENDER_PATH), post.Slug+".html")

	return c.Attachment(fp, post.Slug+".html")
}

func (fe DashboardFrontend) ExportMarkdown(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))

	post, err := fe.PostRepo.GetByID(ctx, uint(id))
	if err != nil {
		return err
	}

	fp := filepath.Join(config.Get(config.POST_RENDER_PATH), "markdowns", post.Slug+".md")

	return c.Attachment(fp, post.Slug+".md")
}

func (fe DashboardFrontend) ExportJSON(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))

	post, err := fe.PostRepo.GetByID(ctx, uint(id))
	if err != nil {
		return err
	}

	// check if temp folder exists
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.Mkdir("temp", 0o755)
	}

	fp := filepath.Join("temp", post.Slug+".json")
	js := []byte(post.Content)
	f, err := os.Create(fp)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(js)
	if err != nil {
		return err
	}

	return c.Attachment(fp, post.Slug+".json")
}
