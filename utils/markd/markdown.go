package markd

import (
	"bytes"
	"fmt"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/dlclark/regexp2"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func ParseMD(source string) (string, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithStyle("xcode-dark"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		panic(err)
	}

	result := buf.String()

	styledResult := changeTagWithClasses(result)

	return styledResult, nil
}

func changeTagWithClasses(str string) string {
	str = strings.ReplaceAll(str,
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
	regex := regexp2.MustCompile("<code>(?!<)", regexp2.None)
	str, _ = regex.Replace(str, `<code class="inline-code">`, 0, -1)
	str = fmt.Sprintf(`<div class="text-justify"> %s </div>`, str)
	return str
}
