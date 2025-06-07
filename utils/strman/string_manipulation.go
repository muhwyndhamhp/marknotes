package strman

import (
	"regexp"
	"strings"
	"unicode"
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

func ToTitleCase(s string) string {
	return strings.Join(strings.FieldsFunc(s, unicode.IsSpace), " ")
}

func ProperTitle(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			r := []rune(word)
			r[0] = unicode.ToTitle(r[0])
			for j := 1; j < len(r); j++ {
				r[j] = unicode.ToLower(r[j])
			}
			words[i] = string(r)
		}
	}
	return strings.Join(words, " ")
}
