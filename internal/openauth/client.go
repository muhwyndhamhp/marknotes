package openauth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/internal"
	"github.com/toolbeam/openauth/client"
	"github.com/toolbeam/openauth/subject"
	"gorm.io/gorm"
)

type authClient struct {
	auth *client.Client
}

func (c *authClient) Callback() func(c echo.Context) error {
	return func(ctx echo.Context) error {
		origin := ctx.Request().URL.String()
		if origin == "" {
			return errors.New("invalid origin")
		}

		baseUrl, err := url.Parse(config.Get(config.BASE_URL))
		if err != nil {
			return err
		}

		redirectURI := baseUrl.
			JoinPath("openauth").
			JoinPath("callback").
			String()

		code := ctx.QueryParam("code")
		if code == "" {
			return errors.New("code is required")
		}

		exchanged, err := c.auth.Exchange(code, redirectURI, &client.ExchangeOptions{})
		if err != nil {
			return err
		}

		if err := c.setCookie(ctx, exchanged.Tokens.Access, exchanged.Tokens.Refresh); err != nil {
			return err
		}

		fmt.Printf("Set-Cookie: %+v", ctx.Response().Header().Get("Set-Cookie"))

		return ctx.Redirect(http.StatusSeeOther, "/dashboard")
	}
}

func (c *authClient) Authorize() (string, error) {
	baseUrl, err := url.Parse(config.Get(config.BASE_URL))
	if err != nil {
		return "", err
	}

	redirectURI := baseUrl.
		JoinPath("openauth").
		JoinPath("callback").
		String()

	authorize, err := c.auth.Authorize(redirectURI, "code", &client.AuthorizeOptions{})
	if err != nil {
		return "", err
	}

	return authorize.URL, nil
}

func NewClient(app *internal.Application) internal.OpenAuth {
	baseOpenAuth, err := url.Parse(config.Get(config.OPEN_AUTH_URL))
	if err != nil {
		panic(err)
	}

	cl, err := client.NewClient(client.ClientInput{
		ClientID: "marknotes-go",
		Issuer:   baseOpenAuth.String(),
		SubjectSchema: subject.SubjectSchemas{
			"user": func(properties any) (any, error) {
				user, ok := properties.(map[string]any)
				if !ok {
					return nil, errors.New("invalid user type")
				}

				if user["userID"] == nil {
					return nil, errors.New("id is required")
				}

				if user["userID"] == "INVALID" {
					return nil, errors.New("user is invalid")
				}

				if user["email"] == nil {
					return nil, errors.New("email is required")
				}

				oauthID := ""

				if user["oauthID"] != nil {
					oauthID = user["oauthID"].(string)
				}

				userId, err := strconv.ParseUint(user["userID"].(string), 10, 64)
				if err != nil {
					return nil, err
				}

				return internal.User{
					Model:       gorm.Model{ID: uint(userId)},
					Email:       user["email"].(string),
					OauthUserID: oauthID,
				}, nil
			},
		},
	})
	if err != nil {
		panic(err)
	}

	return &authClient{
		auth: cl,
	}
}

func (c *authClient) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			cookie, err := ctx.Cookie("access_token")
			if err != nil {
				log.Error(err)

				return ctx.Redirect(http.StatusFound, "/openauth/authorize")
			}

			refreshCookie, err := ctx.Cookie("refresh_token")
			if err != nil {
				log.Error(err)

				return ctx.Redirect(http.StatusFound, "/openauth/authorize")
			}

			verified, err := c.auth.Verify(cookie.Value, &client.VerifyOptions{Refresh: refreshCookie.Value})
			if err != nil {
				log.Error(err)

				return ctx.Redirect(http.StatusFound, "/openauth/authorize")
			}

			if verified.Tokens != nil {
				if err := c.setCookie(ctx, verified.Tokens.Access, verified.Tokens.Refresh); err != nil {
					log.Error(err)

					return ctx.Redirect(http.StatusFound, "/openauth/authorize")
				}
			}

			_, ok := verified.Subject.Properties.(internal.User)

			if !ok {
				log.Errorf("failed to parse user data")

				return ctx.Redirect(http.StatusFound, "/openauth/authorize")
			}

			return next(ctx)
		}
	}
}

func (c *authClient) GetUserFromSession(ctx echo.Context) (*internal.User, error) {
	cookie, err := ctx.Cookie("access_token")
	if err != nil {
		return nil, err
	}

	verified, err := c.auth.Verify(cookie.Value, &client.VerifyOptions{})
	if err != nil {
		return nil, err
	}

	if verified.Tokens != nil {
		if err := c.setCookie(ctx, verified.Tokens.Access, verified.Tokens.Refresh); err != nil {
			return nil, err
		}
	}

	user, ok := verified.Subject.Properties.(internal.User)

	if !ok {
		return nil, err
	}

	return &user, nil
}

func (c *authClient) setCookie(ctx echo.Context, accessToken, refreshToken string) error {
	ctx.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		MaxAge:   34560000,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	})

	ctx.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   34560000,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	})

	return nil
}
