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
		{
			AnchorUrl: "/login",
			Label:     "Login",
			AuthRule:  pub_variables.GuestMode,
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
		},
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
		{
			AnchorUrl: "/dashboard",
			Label:     "Dashboard",
			AuthRule:  pub_variables.UserMode,
			UserID:    userID,
		},
		// {
		// 	AnchorUrl: "/posts/new",
		// 	Label:     "Create Post",
		// 	AuthRule:  pub_variables.UserMode,
		// 	UserID:    userID,
		// },
		// {
		// 	AnchorUrl: "/posts_manage",
		// 	Label:     "Manage Post",
		// 	AuthRule:  pub_variables.UserMode,
		// 	UserID:    userID,
		// },
		{
			AnchorUrl: "/logout",
			Label:     "Logout",
			AuthRule:  pub_variables.UserMode,
			UserID:    userID,
		},
	}
}
