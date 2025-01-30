package internal

import "context"

type RenderFile interface {
	RenderPost(ctx context.Context, post *Post)
	RenderPosts(ctx context.Context)
	RenderMarkdowns(ctx context.Context)
	RenderMarkdown(post *Post) error
}
