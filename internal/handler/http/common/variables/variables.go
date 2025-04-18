package variables

import "github.com/a-h/templ"

type BodyOpts struct {
	HeaderButtons []InlineButton
	FooterButtons []InlineButton
	Component     templ.Component
	ExtraHeaders  []templ.Component
	HideTitle     bool
}

type AuthRule string

type InlineButton struct {
	AnchorUrl string
	Label     string
	AuthRule  AuthRule
	UserID    uint
	IsBoosted bool
	BaseURL   string
}

const (
	HeaderButtonsKey = "HeaderButtons"
	FooterButtonsKey = "FooterButtons"
)

const (
	UserMode   AuthRule = "user-mode"
	GuestMode  AuthRule = "guest-mode"
	AlwaysMode AuthRule = "always-mode"
)

func IsVisible(authRule AuthRule, userID uint) bool {
	if authRule == "always-mode" ||
		(authRule == "user-mode" && userID != 0) ||
		(authRule == "guest-mode" && userID == 0) {
		return true
	} else {
		return false
	}
}
