package pub_variables

import "github.com/a-h/templ"

type DashboardOpts struct {
	Nav  []DrawerMenu
	Comp templ.Component
}

type DrawerMenu struct {
	Label    string
	URL      templ.SafeURL
	IsActive bool
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