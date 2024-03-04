package admin

import pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"

type AuthNeed string

func AppendFooterButtons(userID uint) []pub_variables.InlineButton {
	return []pub_variables.InlineButton{
		{
			AnchorUrl: "/contact",
			Label:     "Contact Me",
			AuthRule:  pub_variables.AlwaysMode,
			UserID:    userID,
		},
		{
			AnchorUrl: "/resume",
			Label:     "Resume",
			AuthRule:  pub_variables.AlwaysMode,
			UserID:    userID,
		},
	}
}

func AppendHeaderButtons(userID uint) []pub_variables.InlineButton {
	return []pub_variables.InlineButton{
		{
			AnchorUrl: "/articles",
			Label:     "Articles",
			AuthRule:  pub_variables.AlwaysMode,
			UserID:    userID,
			IsBoosted: true,
		},
		{
			AnchorUrl: "/contact",
			Label:     "Contact Me",
			AuthRule:  pub_variables.AlwaysMode,
			UserID:    userID,
			IsBoosted: true,
		},
		{
			AnchorUrl: "/resume",
			Label:     "Resume",
			AuthRule:  pub_variables.AlwaysMode,
			UserID:    userID,
			IsBoosted: true,
		},
		{
			AnchorUrl: "/dashboard",
			Label:     "Dashboard",
			AuthRule:  pub_variables.AlwaysMode,
			UserID:    userID,
			IsBoosted: true,
		},
	}
}
