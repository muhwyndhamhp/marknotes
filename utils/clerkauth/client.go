package clerkauth

import (
	"errors"
	"net/http"

	"github.com/apsystole/log"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
)

type Client struct {
	Clerk clerk.Client
}

var sessionCache = map[string]bool{}

func NewClient(secret string) *Client {
	cl, err := clerk.NewClient(secret)
	if err != nil {
		panic(err)
	}

	return &Client{Clerk: cl}
}

// echo middleware to bounce unautorized access to login page
func (cl *Client) AuthMiddleware(userRepo models.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, ok := clerk.SessionFromContext(c.Request().Context())
			if !ok {
				for k := range sessionCache {
					delete(sessionCache, k)
				}

				return c.Redirect(http.StatusFound, "/dashboard/login")
			}

			if valid, ok := sessionCache[session.SessionID]; ok && valid {
				return next(c)
			}

			u, err := cl.GetUser(session)
			if err != nil {
				return c.Redirect(http.StatusFound, "/dashboard/login")
			}

			usr := userRepo.GetCache(c.Request().Context(), u.EmailAddresses[0].EmailAddress)
			if usr == nil {
				_, err = cl.Clerk.Sessions().Revoke(session.SessionID)
				if err != nil {
					log.Error(err)
				}
				return c.Redirect(http.StatusFound, "/dashboard/login")
			}

			sessionCache[session.SessionID] = true

			return next(c)
		}
	}
}

func (cl *Client) GetUserFromSession(c echo.Context, userRepo models.UserRepository) (*models.User, error) {
	session, ok := clerk.SessionFromContext(c.Request().Context())
	if !ok {
		return nil, errors.New("session not found")
	}

	u, err := cl.GetUser(session)
	if err != nil {
		return nil, err
	}

	usr := userRepo.GetCache(c.Request().Context(), u.EmailAddresses[0].EmailAddress)
	if usr == nil {
		return nil, errors.New("user not found")
	}

	return usr, nil
}

func (cl *Client) GetUser(sc *clerk.SessionClaims) (*clerk.User, error) {
	s, err := cl.Clerk.Sessions().Read(sc.SessionID)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	u, err := cl.Clerk.Users().Read(s.UserID)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return u, nil
}
