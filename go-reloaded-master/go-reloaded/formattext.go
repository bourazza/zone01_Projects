package relod

import (
	"strings"
	"unicode"
)

func FormatText(text string) string {
	text = strings.ReplaceAll(text, " ...", "...")
	text = strings.ReplaceAll(text, " !?", "!?")
	text = strings.ReplaceAll(text, " !", "!")
	text = strings.ReplaceAll(text, " ?", "?")
	text = strings.ReplaceAll(text, " .", ".")
	text = strings.ReplaceAll(text, " ,", ",")
	text = strings.ReplaceAll(text, " :", ":")
	text = strings.ReplaceAll(text, " ;", ";")

	result := ""
	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		result += string(runes[i])
		if i < len(runes)-1 && isPunctuation(runes[i]) && !isPunctuation(runes[i+1]) && runes[i+1] != ' ' {
			result += " "
		}
	}

	words := customSplit(result)
	for i := 0; i < len(words)-1; i++ {

		if strings.ToLower(words[i]) == "a" {
			nextWord := words[i+1]
			if len(nextWord) > 0 && isVowelOrH(nextWord[0]) {
				words[i] = "an"
			}
		}
	}
	r := strings.Join(words, " ")

	return r
}
func isPunctuation(r rune) bool {
	return r == '.' || r == ',' || r == '!' || r == '?' || r == ':' || r == ';'
}
func isVowelOrH(c byte) bool {
	return c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' || c == 'h' ||
		c == 'A' || c == 'E' || c == 'I' || c == 'O' || c == 'U' || c == 'H'
}
func customSplit(s string) []string {
	var result []string
	var currentWord string
	inParentheses := false

	for i, r := range s {
		if r == '(' {
			inParentheses = true
			currentWord += string(r)
			if i+1 < len(s) && rune(s[i+1]) != ')' {
				continue
			}
		} else if r == ')' {
			inParentheses = false
			currentWord += string(r)
		} else if unicode.IsSpace(r) && !inParentheses {
			if currentWord != "" {
				result = append(result, currentWord)
				currentWord = ""
			}
		} else {
			currentWord += string(r)
		}

		if i == len(s)-1 && currentWord != "" {
			result = append(result, currentWord)
		}
	}

	return result
}
