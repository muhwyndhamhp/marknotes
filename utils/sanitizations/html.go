package sanitizations

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

func SanitizeHtml(escapedHTML string) string {
	p := bluemonday.UGCPolicy()

	p.AllowAttrs("class").
		Matching(
			regexp.MustCompile(
				`(hljs-string)`+`|(hljs-title)`+
					`|(hljs-function)`+`|(hljs-params)`+
					`|(hljs-keyword)`+`|(hljs-type)`+
					`|(hljs-number)`+`|(hljs-comment)`+
					`|(\bsuggestion\b)`+`|(\blanguage-([a-z]*\b)\b)`,
			)).
		OnElements("span", "code")

	p.AllowAttrs("contenteditable").Matching(regexp.MustCompile(`(false)`)).OnElements("span")
	p.AllowAttrs("data-type").Matching(regexp.MustCompile(`(mention)`)).OnElements("span")
	p.AllowAttrs("data-id").Matching(regexp.MustCompile(`([a-z]+)`)).OnElements("span")

	html := p.Sanitize(escapedHTML)

	return html
}
