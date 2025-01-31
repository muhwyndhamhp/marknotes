package admin

import (
	"strings"

	"github.com/muhwyndhamhp/marknotes/config"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
)

type AuthNeed string

func AppendHeaderButtons(userID uint) []pub_variables.InlineButton {
	baseURL := strings.Split(config.Get(config.OAUTH_URL), "/callback")[0]

	return []pub_variables.InlineButton{
		{
			AnchorUrl: baseURL + "/articles",
			Label:     "Articles",
			AuthRule:  pub_variables.AlwaysMode,
			UserID:    userID,
			IsBoosted: true,
			BaseURL:   baseURL,
		},
		{
			AnchorUrl: baseURL + "/resume",
			Label:     "Resume",
			AuthRule:  pub_variables.AlwaysMode,
			UserID:    userID,
			IsBoosted: true,
			BaseURL:   baseURL,
		},
		{
			AnchorUrl: baseURL + "/dashboard",
			Label:     "Dashboard",
			AuthRule:  pub_variables.AlwaysMode,
			UserID:    userID,
			IsBoosted: true,
			BaseURL:   baseURL,
		},
	}
}
