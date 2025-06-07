package internal

import "context"

type ProfanityCheck interface {
	IsProfane(ctx context.Context, text string) bool
}
