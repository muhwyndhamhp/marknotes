package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/internal"
	"github.com/samber/lo"
	"time"
)
import "google.golang.org/genai"

type client struct {
	app       *internal.Application
	genClient *genai.Client
}

func (c client) ModerateReplies(ctx context.Context, replies []internal.Reply) ([]internal.Reply, error) {
	cfg := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"id":                {Type: genai.TypeInteger},
					"message":           {Type: genai.TypeString},
					"moderation_status": {Type: genai.TypeInteger},
					"moderation_reason": {Type: genai.TypeString},
				},
				PropertyOrdering: []string{"id", "message", "moderation_status", "moderation_reason"},
			},
		},
	}

	req := lo.Reduce(replies, func(agg string, item internal.Reply, index int) string {
		i := map[string]interface{}{
			"id":      item.ID,
			"message": item.Message,
		}

		js, _ := json.Marshal(i)

		if index == 0 {
			agg = fmt.Sprintf("[%s", string(js))
		} else {
			agg = fmt.Sprintf("%s,%s", agg, string(js))
		}

		if index == len(replies)-1 {
			agg = agg + "]"
		}

		return agg
	}, "")

	fmt.Println(req)

	result, err := c.genClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-preview-05-20",
		genai.Text(fmt.Sprintf("%s \n %s", moderateReplyPrompt, req)),
		cfg,
	)
	if err != nil {
		return nil, err
	}

	var mods []struct {
		ID               uint   `json:"id"`
		Message          string `json:"message"`
		ModerationStatus uint   `json:"moderation_status"`
		ModerationReason string `json:"moderation_reason"`
	}

	err = json.Unmarshal([]byte(result.Text()), &mods)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	replies = lo.Map(replies, func(item internal.Reply, index int) internal.Reply {
		m, ok := lo.Find(mods, func(item2 struct {
			ID               uint   `json:"id"`
			Message          string `json:"message"`
			ModerationStatus uint   `json:"moderation_status"`
			ModerationReason string `json:"moderation_reason"`
		}) bool {
			return item.ID == item2.ID
		})

		if !ok {
			return internal.Reply{}
		}

		item.Moderation = internal.Moderation{
			LastModeratedAt:  &now,
			ModerationStatus: internal.ModerationStatus(m.ModerationStatus),
			ModerationReason: m.ModerationReason,
		}

		return item
	})

	return replies, nil
}

const moderateReplyPrompt = `
Imagine you are content moderator for a blog site revolving around technology and programming. 
Your task is to moderate following comments.

'id' is an integer and should return the same id that supplied from the input.

'message' is a string and should return the same message that supplied from the input.

'moderation_status' is an integer, where 1 means "OK", 2 means "Warning", and 3 means "Dangerous".
- OK means all comments that uses normal language, including disagreement and civil debates.
- Warning means it may have strong language, implicit mockery, swearing in context, and any normal message that contains URL.
- Dangerous should cover actual hate speech, slander, slur, ad hominem, straw man fallacy, and spam / unrelated content promotion including placeholder / lipsum texts.

'moderation_reason' is a string that contains single sentence, no more than 15 words summarizing the 
reasoning for the moderation status. 

Comments will be provided in a JSON format of { "id": integer, "message": string }, and will be attached
right after this message.
`

func NewClient(app *internal.Application) internal.LLM {
	c, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  config.Get(config.GEMINI_API_KEY),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		panic(err)
	}

	return &client{app: app, genClient: c}
}
