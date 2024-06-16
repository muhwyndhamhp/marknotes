package renderfile

import (
	"context"
	"errors"
	"fmt"
	"github.com/muhwyndhamhp/marknotes/db"
	"os"

	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"github.com/muhwyndhamhp/marknotes/utils/fileman"
)

func RenderMarkdowns(ctx context.Context, repo models.PostRepository) {
	if err := fileman.CheckDir(config.Get(config.POST_RENDER_PATH) + "/markdowns"); err != nil {
		fmt.Println(err)
	}
	posts, err := repo.Get(ctx, db.Where("status = ?", values.Published))
	if err != nil {
		fmt.Println(err)
	}

	for _, post := range posts {
		err := RenderMarkdown(&post)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func RenderMarkdown(post *models.Post) error {
	if post.MarkdownContent == "" {
		return errors.New("post content is empty")
	}

	err := os.WriteFile(
		config.Get(config.POST_RENDER_PATH)+"/markdowns/"+post.Slug+".md",
		[]byte(post.MarkdownContent),
		0o755,
	)
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}
