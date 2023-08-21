package strman

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
