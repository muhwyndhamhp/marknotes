package profanity

import (
	"context"
	"unicode"

	goaway "github.com/TwiN/go-away"
	"github.com/muhwyndhamhp/marknotes/internal"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type client struct {
	detector *goaway.ProfanityDetector
}

func NewClient() internal.ProfanityCheck {
	profanityDetector := goaway.NewProfanityDetector().WithSanitizeSpaces(true)

	return &client{detector: profanityDetector}
}

func (c *client) IsProfane(ctx context.Context, text string) bool {
	return c.detector.IsProfane(c.removeNonASCII(text))
}

func (c *client) removeNonASCII(input string) string {
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	result, _, _ := transform.String(t, input)
	return c.sanitizeASCII(result)
}

func (c *client) sanitizeASCII(input string) string {
	var output []rune
	for _, r := range input {
		if r <= 127 { // ASCII range
			output = append(output, r)
		}
	}
	return string(output)
}
