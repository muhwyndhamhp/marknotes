package admin

type InlineButton struct {
	AnchorUrl string
	Label     string
	AuthRule  AuthRule
	UserID    uint
}

type AuthRule string

const (
	UserMode   AuthRule = "user-mode"
	GuestMode  AuthRule = "guest-mode"
	AlwaysMode AuthRule = "always-mode"

	HeaderButtonsKey = "HeaderButtons"
	FooterButtonsKey = "FooterButtons"
)

type AuthNeed string

func AppendFooterButtons(userID uint) []InlineButton {
	return []InlineButton{
		{
			AnchorUrl: "/contact",
			Label:     "Contact Me",
			AuthRule:  AlwaysMode,
			UserID:    userID,
		},
		{
			AnchorUrl: "/resume",
			Label:     "Resume",
			AuthRule:  AlwaysMode,
			UserID:    userID,
		},
		{
			AnchorUrl: "/login",
			Label:     "Login",
			AuthRule:  GuestMode,
			UserID:    userID,
		},
	}
}

func AppendHeaderButtons(userID uint) []InlineButton {
	return []InlineButton{
		{
			AnchorUrl: "/posts_index",
			Label:     "Articles",
			AuthRule:  AlwaysMode,
			UserID:    userID,
		},
		{
			AnchorUrl: "/contact",
			Label:     "Contact Me",
			AuthRule:  AlwaysMode,
			UserID:    userID,
		},
		{
			AnchorUrl: "/resume",
			Label:     "Resume",
			AuthRule:  AlwaysMode,
			UserID:    userID,
		},
		{
			AnchorUrl: "/posts/new",
			Label:     "Create Post",
			AuthRule:  UserMode,
			UserID:    userID,
		},
		{
			AnchorUrl: "/posts_manage",
			Label:     "Manage Post",
			AuthRule:  UserMode,
			UserID:    userID,
		},
		{
			AnchorUrl: "/logout",
			Label:     "Logout",
			AuthRule:  UserMode,
			UserID:    userID,
		},
	}
}
