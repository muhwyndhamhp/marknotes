package admin

import (
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/common/variables"
	"strings"

	"github.com/muhwyndhamhp/marknotes/config"
)

type AuthNeed string

func AppendHeaderButtons(userID uint) []variables.InlineButton {
	baseURL := strings.Split(config.Get(config.OAUTH_URL), "/callback")[0]

	return []variables.InlineButton{
		{
			AnchorUrl: baseURL + "/articles",
			Label:     "Articles",
			AuthRule:  variables.AlwaysMode,
			UserID:    userID,
			IsBoosted: true,
			BaseURL:   baseURL,
		},
		{
			AnchorUrl: baseURL + "/resume",
			Label:     "Resume",
			AuthRule:  variables.AlwaysMode,
			UserID:    userID,
			IsBoosted: true,
			BaseURL:   baseURL,
		},
		{
			AnchorUrl: baseURL + "/dashboard",
			Label:     "Dashboard",
			AuthRule:  variables.AlwaysMode,
			UserID:    userID,
			IsBoosted: true,
			BaseURL:   baseURL,
		},
	}
}
