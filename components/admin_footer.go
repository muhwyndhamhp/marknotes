package components

import "github.com/muhwyndhamhp/marknotes/components/elements"

type AdminFooter struct {
	Buttons      []elements.InlineAnchorButton
	AdminButtons []elements.InlineAnchorButton
	UserID       uint
}

func GetAdminFooter(userID uint) AdminHeader {
	return AdminHeader{
		Buttons: []elements.InlineAnchorButton{
			{
				AnchorLink: "/contact",
				Label:      "Contact Me",
			},
			{
				AnchorLink: "/resume",
				Label:      "Resume",
			},
		},
		AdminButtons: []elements.InlineAnchorButton{
			{
				AnchorLink: "/login",
				Label:      "Login",
			},
		},
		UserID: userID,
	}
}
