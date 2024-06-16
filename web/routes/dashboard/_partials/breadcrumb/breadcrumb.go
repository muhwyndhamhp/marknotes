package breadcrumb

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/muhwyndhamhp/marknotes/config"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func Breadcrumbs(path string) []pub_variables.Breadcrumb {
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
