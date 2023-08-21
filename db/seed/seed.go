package main

import (
	"context"
	"html/template"

	"github.com/muhwyndhamhp/marknotes/db"
	"github.com/muhwyndhamhp/marknotes/pkg/repository"
	"github.com/muhwyndhamhp/marknotes/utils/markd"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"github.com/muhwyndhamhp/marknotes/utils/strman"
)

func main() {

	ctx := context.Background()
	repo := repository.NewPostRepository(db.GetDB().Debug())

	posts, _ := repo.Get(ctx, scopes.QueryOpts{})
	for i := range posts {
		str := strman.TakeFirstWords(40, posts[i].Content)
		str = strman.AddTrailingComma(str)
		preview, _ := markd.ParseMD(str)
		posts[i].Preview = template.HTML(preview)
		repo.Upsert(ctx, &posts[i])
	}
}
