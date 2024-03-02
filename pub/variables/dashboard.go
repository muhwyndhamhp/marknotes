package pub_variables

import "github.com/a-h/templ"

type DashboardOpts struct {
	Nav               []DrawerMenu
   BreadCrumbs       []Breadcrumb
	Comp              templ.Component
	AdditionalHeaders []string
	NavItems          []templ.Component
}

type DrawerMenu struct {
	Label     string
	URL       templ.SafeURL
	IsActive  bool
	IsBoosted bool
}

type Breadcrumb struct {
	Label string
	URL   templ.SafeURL
}

type DropdownVM struct {
	Items    []DropdownItem
	Selected int
}

type DropdownItem struct {
	Label  string
	Path   string
	Target string
}
