package markd

import (
	"bytes"
	"strings"

	embed "github.com/13rac1/goldmark-embed"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var mdParser goldmark.Markdown

func init() {
	mdParser = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithStyle("xcode-dark"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
			embed.New(),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)
}

func ParseMD(source string) (string, error) {
	var buf bytes.Buffer
	if err := mdParser.Convert([]byte(source), &buf); err != nil {
		return "", err
	}

	result := buf.String()
	styledResult := postProcessHTML(result)

	return styledResult, nil
}

func postProcessHTML(str string) string {
	str = changeCodeBlockBg(str)
	return str
}

func changeCodeBlockBg(str string) string {
	return strings.ReplaceAll(str,
		`style="color:#fff;background-color:#1f1f24;"`,

		`style="color:#fff;
		background-color:rgb(15 23 42);
		width:100%;
		border-radius:0.5rem;
		margin-top: 1.5rem;
		margin-bottom: 1.5rem;
		padding: 1.25rem;
		overflow-x: scroll;
		"
		`)
}
