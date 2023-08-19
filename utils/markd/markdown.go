package markd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	HEADER_1 = `<h1 class="text-5xl text-teal-600 my-2" `
	HEADER_2 = `<h2 class="text-3xl text-teal-600 my-2" `
	HEADER_3 = `<h3 class="text-2xl text-teal-600 my-2" `
	HEADER_4 = `<h4 class="text-xl text-teal-600 my-2" `
	HEADER_5 = `<h5 class="text-lg text-teal-600 my-2" `
)

func ParseMD(source string) (string, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
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

	str = strings.ReplaceAll(str, "<h1 ", HEADER_1)
	str = strings.ReplaceAll(str, "<h2 ", HEADER_2)
	str = strings.ReplaceAll(str, "<h3 ", HEADER_3)
	str = strings.ReplaceAll(str, "<h4 ", HEADER_4)
	str = strings.ReplaceAll(str, "<h5 ", HEADER_5)

	str = fmt.Sprintf(`<div class="text-justify"> %s </div>`, str)
	return str
}
