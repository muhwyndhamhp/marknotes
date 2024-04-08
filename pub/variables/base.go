package pub_variables

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
