package sanitizations

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/muhwyndhamhp/marknotes/config"
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
					`|(hljs-variable)`+`|(hljs-selector-class)`+
					`|(hljs-selector-id)`+`|(hljs-selector-tag)`+
					`|(hljs-meta)`+`|(hljs-tag)`+`|(hljs-attribute)`+
					`|(\bsuggestion\b)`+`|(\blanguage-([a-z]*\b)\b)`+
					`|(\bxml\b)`+`|(\bmockup-code\b)`,
			)).
		OnElements("span", "code", "pre")

	p.AllowAttrs("loading").OnElements("img")
	p.AllowAttrs("hx-get").OnElements("div")
	p.AllowAttrs("hx-swap").OnElements("div")
	p.AllowAttrs("hx-trigger").OnElements("div")
	p.AllowAttrs("contenteditable").Matching(regexp.MustCompile(`(false)`)).OnElements("span")
	p.AllowAttrs("data-type").Matching(regexp.MustCompile(`(mention)`)).OnElements("span")
	p.AllowAttrs("data-id").Matching(regexp.MustCompile(`([a-z]+)`)).OnElements("span")

	p.AllowStandardURLs()
	p.AllowAttrs("class").Matching(regexp.MustCompile(`(max-h-96)|(mx-auto)`)).OnElements("img")
	p.AllowURLSchemes("mailto", "http", "https")
	p.AllowAttrs("src").OnElements("img")
	p.RequireParseableURLs(true)

	p.AllowAttrs("class").Matching(regexp.MustCompile(`(mx-auto)`)).OnElements("iframe")
	p.AllowAttrs("width", "height").Matching(regexp.MustCompile(`([0-9]+)`)).OnElements("iframe")

	if config.Get(config.ENV) == "dev" {
		p.AllowRelativeURLs(true)
	}

	html := p.Sanitize(escapedHTML)

	return html
}
