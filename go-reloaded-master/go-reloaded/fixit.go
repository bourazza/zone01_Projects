package relod

import (
	"fmt"
	"strings"
)

func Fixit(text string) string {
	quoteCount := strings.Count(text, "'")
	result := ""
	var k int
	if quoteCount%2 != 0 {
		k = quoteCount - 1
	} else {
		k = quoteCount
	}
	first := true
	rm_next := false
	if len(text) <= 3 {
		return text
	}
	for idx, char := range text {
		if char == '\'' && k != 0 {
			fmt.Print("kkk ")
			if idx == 0 {
				result += string(char)
				first = false
				if text[idx+1] == ' ' {
					rm_next = true
				}
			} else if idx == len(text)-1 {
				fmt.Print("h ")
				if text[idx-1] == ' ' {
					fmt.Print("h ")
					result = result[:len(result)-1]
				}
				result += "'"
			} else {
				if text[idx-1] != ' ' && text[idx+1] != ' ' {
					result += "'"
					continue
				}
				if first {
					if result[len(result)-1] != ' ' && result[len(result)-1] != '\'' {
						result += " "
					}
					result += "'"
					if text[idx+1] == ' ' {
						rm_next = true
					}
					first = false
				} else {
					first = true
					if result[len(result)-1] == ' ' {
						result = result[:len(result)-1]
					}
					is_punc := strings.Contains(",;:.!? ", string(text[idx+1]))
					result += "'"
					if !is_punc {
						result += " "
					}
				}
			}
			k = k - 1
		} else {
			if rm_next {
				rm_next = false
				continue
			}
			result += string(char)
		}
	}

	return result
}
