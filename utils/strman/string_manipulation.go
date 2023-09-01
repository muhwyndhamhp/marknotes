package strman

import (
	"regexp"
	"strings"
)

const (
	regexMatchWords = `\b[A-Za-z]+\b`
)

func TakeFirstWords(wordCount int, source string) string {
	count := 0

	for i := range source {
		if source[i] != ' ' {
			continue
		}

		count++
		if count >= wordCount {
			return source[:i]
		}
	}

	return source
}

func AddTrailingComma(str string) string {
	return str + "..."
}

func GenerateSlug(source string) (string, error) {
	re := regexp.MustCompile(regexMatchWords)

	cleanTitle := re.FindAllString(source, -1)

	str := strings.ToLower(TakeFirstWords(8, strings.Join(cleanTitle[:], " ")))
	return strings.ReplaceAll(str, " ", "-"), nil
}
