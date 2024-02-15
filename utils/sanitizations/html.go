package sanitizations

import "github.com/microcosm-cc/bluemonday"

func SanitizeHtml(escapedHTML string) string {
	p := bluemonday.UGCPolicy()

	html := p.Sanitize(escapedHTML)

	return html
}
