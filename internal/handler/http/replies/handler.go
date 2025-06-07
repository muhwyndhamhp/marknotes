package replies

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/internal"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/replies/reply"
	"github.com/muhwyndhamhp/marknotes/internal/middlewares"
	templateRenderer "github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/strman"
	"github.com/sethvargo/go-diceware/diceware"
	"net/http"
)

type handler struct {
	app *internal.Application
}

func NewHandler(g *echo.Group, app *internal.Application) {
	h := &handler{app: app}

	g.GET("/replies/articles/:id", h.ArticleReplies, app.GetIdParamWare)
	g.POST("/replies/articles/:id/create", h.Create, app.GetIdParamWare)
}

func (h *handler) ArticleReplies(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := c.Get(middlewares.ByIDKey).(int)

	replies, err := h.app.ReplyRepository.FetchArticleReplies(ctx, uint(id))
	if err != nil {
		return err
	}

	replyTemplate := reply.ArticleReplies(uint(id), replies)

	return templateRenderer.AssertRender(c, http.StatusOK, replyTemplate)
}

type CreateReplyReq struct {
	ReplyBody string `json:"reply_body" form:"replyBody"`
	ParentID  uint   `json:"parent_id" form:"parentId"`
}

func (h *handler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := c.Get(middlewares.ByIDKey).(int)

	var req CreateReplyReq

	if err := c.Bind(&req); err != nil {
		return err
	}

	if hasProfanity := h.app.ProfanityCheck.IsProfane(ctx, req.ReplyBody); hasProfanity {
		// TODO add proper error message to text input on response
		return errors.New("input has profanity in it")
	}

	existingAlias, _ := c.Request().Cookie("comment_alias")
	list, _ := diceware.Generate(2)
	alias := strman.ProperTitle(fmt.Sprintf("%s %s", list[0], list[1]))

	if existingAlias != nil {
		alias = existingAlias.Value
	} else {
		c.SetCookie(&http.Cookie{
			Name:     "comment_alias",
			Value:    alias,
			MaxAge:   34560000,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
		})
	}

	r := internal.Reply{
		ArticleID: uint(id),
		Message:   req.ReplyBody,
		Alias:     alias,
	}

	if req.ParentID != 0 {
		r.ParentID = &req.ParentID
	}

	if err := h.app.ReplyRepository.CreateReply(ctx, &r); err != nil {
		return err
	}

	return h.ArticleReplies(c)
}
