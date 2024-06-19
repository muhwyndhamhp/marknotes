package sidebar

import pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"

func Nav(indexSelected int) []pub_variables.DrawerMenu {
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
		//{
		//	Label:     "Logout",
		//	URL:       "/logout",
		//	IsActive:  false,
		//	IsBoosted: false,
		//},
	}

	lists[indexSelected].IsActive = true

	return lists
}
