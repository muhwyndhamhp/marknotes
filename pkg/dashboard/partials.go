package dashboard

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/config"
	pub_editor "github.com/muhwyndhamhp/marknotes/pub/components/editor"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	templates "github.com/muhwyndhamhp/marknotes/template"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (fe *DashboardFrontend) Editor(c echo.Context) error {
	baseURL := strings.Split(config.Get(config.OAUTH_URL), "/callback")[0]
	uploadURL := fmt.Sprintf("%s/posts/%d/media/upload", baseURL, 0)

	dashboard := pub_editor.Editor(uploadURL)

	return templates.AssertRender(c, http.StatusOK, dashboard)
}

func (fe *DashboardFrontend) Breadcrumbs(path string) []pub_variables.Breadcrumb {
	baseURL := strings.Split(config.Get(config.OAUTH_URL), "/callback")[0]

	paths := strings.Split(path, "/")
	var items []pub_variables.Breadcrumb
	c := cases.Title(language.English)
	for i, path := range paths {
		if path != "" {
			items = append(items, pub_variables.Breadcrumb{
				Label: c.String(path),
				URL:   templ.SafeURL(fmt.Sprintf("%s/%s", baseURL, strings.Join(paths[:i+1], "/"))),
			})
		}
	}
	return items
}

func (fe *DashboardFrontend) SizeDropdown(page, pageSize int) pub_variables.DropdownVM {
	arrays := []pub_variables.DropdownItem{}
	item := 0
	for i := range []int{0, 1, 2} {
		size := (i + 1) * 10
		arrays = append(arrays, pub_variables.DropdownItem{
			Label:  fmt.Sprintf("%d", size),
			Path:   fmt.Sprintf("/dashboard/articles?page=%d&pageSize=%d&source=source-partial", page, size),
			Target: "#articles",
		})
		if size == pageSize {
			item = i
		}
	}
	return pub_variables.DropdownVM{
		Items:    arrays,
		Selected: item,
	}
}

func (fe *DashboardFrontend) PageDropdown(page, pageSize, count int) pub_variables.DropdownVM {
	arrays := []pub_variables.DropdownItem{}
	item := 0
	for i := 0; (i)*pageSize <= count; i++ {
		currentPage := i + 1
		size := pageSize
		arrays = append(arrays, pub_variables.DropdownItem{
			Label:  fmt.Sprintf("%d", currentPage),
			Path:   fmt.Sprintf("/dashboard/articles?page=%d&pageSize=%d&source=source-partial", currentPage, size),
			Target: "#articles",
		})
		if currentPage == page {
			item = i
		}
	}

	return pub_variables.DropdownVM{
		Items:    arrays,
		Selected: item,
	}
}

func nav(indexSelected int) []pub_variables.DrawerMenu {
	lists := []pub_variables.DrawerMenu{
		{
			Label:     "Articles",
			URL:       "/dashboard/articles",
			IsActive:  false,
			IsBoosted: true,
		},
		{
			Label:     "Back to Site",
			URL:       "/",
			IsActive:  false,
			IsBoosted: true,
		},
	}

	lists[indexSelected].IsActive = true

	return lists
}
