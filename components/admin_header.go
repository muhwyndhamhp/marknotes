package components

import "github.com/muhwyndhamhp/marknotes/components/elements"

type AdminHeader struct {
	Buttons      []elements.InlineAnchorButton
	AdminButtons []elements.InlineAnchorButton
	UserID       uint
}

func GetAdminHeader(userID uint) AdminHeader {
	return AdminHeader{
		Buttons: []elements.InlineAnchorButton{
			{
				AnchorLink: "/posts_index",
				Label:      "Articles",
			},
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
				AnchorLink: "/posts/new",
				Label:      "Create Post",
			},
			{
				AnchorLink: "/posts_manage",
				Label:      "Manage Posts",
			},
			{
				AnchorLink: "/logout",
				Label:      "Logout",
			},
		},
		UserID: userID,
	}
}
