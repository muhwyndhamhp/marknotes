package renderfile

import (
	"context"
	"errors"
	"fmt"
	"github.com/muhwyndhamhp/marknotes/internal"
	"os"

	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"github.com/muhwyndhamhp/marknotes/utils/fileman"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
)

func (r *RenderClient) RenderMarkdowns(ctx context.Context) {
	if err := fileman.CheckDir(config.Get(config.POST_RENDER_PATH) + "/markdowns"); err != nil {
		fmt.Println(err)
	}
	posts, err := r.App.PostRepository.Get(ctx, scopes.Where("status = ?", internal.PostStatusPublished))
	if err != nil {
		fmt.Println(err)
	}

	for _, post := range posts {
		err := r.RenderMarkdown(&post)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (r *RenderClient) RenderMarkdown(post *internal.Post) error {
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
