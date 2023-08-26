package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/utils/jwt"

	"github.com/gorilla/csrf"
)

const (
	userURL = "https://api.github.com/user"
)

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type AuthService struct {
	JWT            jwt.Service
	AuthURL        string
	AccessTokenURL string
	ClientID       string
	ClientSecret   string
	RedirectURL    string
	Repo           models.UserRepository
}

func NewAuthService(g *echo.Group,
	JWT jwt.Service,
	AuthURL, AccessTokenURL,
	ClientID, ClientSecret,
	RedirectURL string,
	Repo models.UserRepository,
) {
	handler := &AuthService{
		JWT:            JWT,
		AuthURL:        AuthURL,
		AccessTokenURL: AccessTokenURL,
		ClientID:       ClientID,
		ClientSecret:   ClientSecret,
		RedirectURL:    RedirectURL,
		Repo:           Repo,
	}

	g.GET("/login", handler.Login)
	g.GET("/callback", handler.Callback)
	g.GET("/logout", handler.Logout)
}

func (h *AuthService) Logout(c echo.Context) error {
	if err := h.JWT.ForgetToken(c); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/posts_index")
}

func (h *AuthService) Callback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	ctx := c.Request().Context()

	expectedState, err := c.Request().Cookie("csrf_token")
	if err != nil {
		return err
	}
	if err == http.ErrNoCookie {
		return c.JSON(http.StatusUnauthorized, err)
	}

	if state != expectedState.Value {
		return c.JSON(http.StatusUnauthorized, err)
	}

	accessToken, err := h.getAccessToken(code)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	oauthUser, err := h.getUser(accessToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	user, _ := h.Repo.GetByOauthID(ctx, fmt.Sprintf("%d", oauthUser.ID))
	if user == nil {
		return c.Redirect(http.StatusFound, "/unauthorized")
	}

	if err := h.JWT.GenerateTokenAndStore(c, user.ID, user.Name); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/posts_index")
}

func (h *AuthService) Login(c echo.Context) error {
	token := csrf.Token(c.Request())

	params := url.Values{
		"client_id":    []string{h.ClientID},
		"redirect_uri": []string{h.RedirectURL},
		"scope":        []string{"read:user,user:email"},
		"state":        []string{token},
	}

	fmt.Println(params)
	u, err := url.ParseRequestURI(h.AuthURL)
	if err != nil {
		return err
	}
	u.RawQuery = params.Encode()

	cookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Minute),
		HttpOnly: true,
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, u.String())
}

type user struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

func (h *AuthService) getUser(accessToken string) (*user, error) {
	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var u user
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (h *AuthService) getAccessToken(code string) (string, error) {
	params := url.Values{
		"client_id":     []string{h.ClientID},
		"client_secret": []string{h.ClientSecret},
		"code":          []string{code},
	}

	u, err := url.ParseRequestURI(h.AccessTokenURL)
	if err != nil {
		return "", err
	}

	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var accTokenResp accessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&accTokenResp); err != nil {
		return "", err
	}
	return accTokenResp.AccessToken, nil
}
