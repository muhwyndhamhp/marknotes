package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/muhwyndhamhp/marknotes/cmd"
	"github.com/muhwyndhamhp/marknotes/internal"
	llm2 "github.com/muhwyndhamhp/marknotes/internal/llm"
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()

	app := cmd.Bootstrap()
	llm := llm2.NewClient(app)

	replies, _, err := app.ReplyRepository.Fetch(ctx)
	if err != nil {
		panic(err)
	}

	replies = append(replies, internal.Reply{
		Model:   gorm.Model{ID: 9},
		Message: "What the fuck are you talking about?",
	})

	replies = append(replies, internal.Reply{
		Model:   gorm.Model{ID: 10},
		Message: "Cryp70 $$$ H3r3!",
	})

	replies = append(replies, internal.Reply{
		Model:   gorm.Model{ID: 11},
		Message: "Yeah I understand what you mean, dealing with array fucking sucks",
	})

	replies = append(replies, internal.Reply{
		Model:   gorm.Model{ID: 12},
		Message: "This guy is a dumbass",
	})

	replies = append(replies, internal.Reply{
		Model:   gorm.Model{ID: 13},
		Message: "Look at this stupid bitch",
	})

	replies = append(replies, internal.Reply{
		Model:   gorm.Model{ID: 14},
		Message: "I'm not taking advice from people that use Golang like they use Java",
	})

	res, err := llm.ModerateReplies(ctx, replies)
	if err != nil {
		panic(err)
	}

	js, _ := json.MarshalIndent(res, "", "	")

	fmt.Println(string(js))
}
