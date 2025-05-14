package relod

import (
	"fmt"
	"regexp"

	"strconv"
	"strings"
	"unicode"
)

func hexToDecimal(hex string) (string, error) {
	decimal, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", decimal), nil
}

func binToDecimal(bin string) (string, error) {
	decimal, err := strconv.ParseInt(bin, 2, 64)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", decimal), nil
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(strings.ToLower(s))

	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func extractNumber(sr string) int {
	parts := strings.Split(sr, ",")
	if len(parts) != 2 {
		return 1
	}
	numStr := strings.Trim(parts[1], " )")
	if strings.HasPrefix(numStr, "-") {

		return 0
	}
	num, err := strconv.Atoi(numStr)
	if err != nil {

		return -1
	}
	return num
}


func TransformText(text string) string {

	j := FormatFlags(text)

	words := strings.Split(j, " ")

	result := make([]string, len(words))

	copy(result, words)

	for i := 0; i < len(result); i++ {

		if i >= 0 {
			currentWord := result[i]

			if currentWord == "(hex)" {
				if i != 0 {
					if result[i-1] == "" {
						d := 2
						for {

							if len(result) > d && i-d >= 0 {
								dec, err := hexToDecimal(result[i-d])
								if err == nil || result[i-d] != "" {
									if dec != "" {
										result[i-d] = dec
										break
									} else {
										break
									}
								} else {
									d += 1
								}

							} else if i-d <= 0 {
								break
							}

						}
					}
					if dec, err := hexToDecimal(result[i-1]); err == nil {
						if dec != "" {
							result[i-1] = dec
						}

					}

					result[i] = ""

					continue
				} else {
					result[i] = ""

					continue

				}
			}

			if currentWord == "(bin)" {
				if i != 0 {

					if result[i-1] == "" {
						d := 2
						for {

							if len(result) > d && i-d >= 0 {
								de, err := binToDecimal(result[i-d])
								if err == nil || result[i-d] != "" {
									if de != "" {
										result[i-d] = de
									}
									break

								} else {
									d += 1

								}

							} else if i-d <= 0 {
								break
							}

						}
					}

					if dec, err := binToDecimal(result[i-1]); err == nil {
						result[i-1] = dec

					}
					result[i] = ""

					continue
				} else {
					result[i] = ""

					continue

				}
			}

			if strings.HasPrefix(currentWord, "(up") && (strings.HasSuffix(currentWord, ")") && !(strings.HasSuffix(currentWord, "))"))) {
				if currentWord == "(up)" {
					if i != 0 {
						if result[i-1] == "" {
							d := 2
							for {

								if len(result) > d && i-d >= 0 {
									dec := strings.ToUpper(result[i-d])
									if result[i-d] != "" {
										result[i-d] = dec
										break
									} else {
										d += 1
									}

								} else if i-d <= 0 {
									break
								}

							}
						}
						if i-1 > 0 && result[i-1] == "" {
							result[i-2] = strings.ToUpper(result[i-2])
							result[i] = ""
							continue

						} else {
							result[i-1] = strings.ToUpper(result[i-1])
							result[i] = ""
							continue
						}
					} else {
						result[i] = ""
						continue
					}
				} else {
					count := extractNumber(currentWord)
					if count != -1 {
						for j := 0; j < count && i-j > 0; j++ {
							if isit(result[i-j-1]) {
								count++

							}
							result[i-j-1] = strings.ToUpper(result[i-j-1])
						}

						result[i] = ""
						continue
					}
				}
			}

			if strings.HasPrefix(currentWord, "(low") && strings.HasSuffix(currentWord, ")") && !(strings.HasSuffix(currentWord, "))")) {

				if currentWord == "(low)" {
					if i != 0 {
						if result[i-1] == "" {

							d := 2
							for {

								if len(result) > d && i-d >= 0 {
									dec := strings.ToLower(result[i-d])
									if result[i-d] != "" {
										result[i-d] = dec
										break
									} else {
										d += 1
									}

								} else if i-d <= 0 {
									break
								}

							}

						} else {

							result[i-1] = strings.ToLower(result[i-1])

						}
					}
					result[i] = ""
					continue

				} else {
					count := extractNumber(currentWord)
					if count != -1 {
						for j := 0; j < count && i-j > 0; j++ {
							if isit(result[i-j-1]) {
								count++

							}
							result[i-j-1] = strings.ToLower(result[i-j-1])
						}

						result[i] = ""
						continue
					}
				}

			}

			if strings.HasPrefix(currentWord, "(cap") && (strings.HasSuffix(currentWord, ")"))&& !(strings.HasSuffix(currentWord, "))")) {

				if currentWord == "(cap)" {
					if i != 0 {
						if result[i-1] == "" {
							d := 2
							for {

								if len(result) > d && i-d >= 0 {
									dec := capitalize(result[i-d])
									if result[i-d] != "" {
										result[i-d] = dec
										break
									} else {
										d += 1
									}

								} else if i-d <= 0 {
									break
								}

							}
						}

						if  i-1 > 0 && result[i-1] == "" {
							result[i-2] = capitalize(result[i-2])
							result[i] = ""
							continue

						} else {
							result[i-1] = capitalize(result[i-1])
							result[i] = ""
							continue
						}
					} else {
						result[i] = ""
						continue
					}
				} else {

					count := extractNumber(currentWord)
					if count != -1 {
						

						if count > i {
							count = i
						}
						for j := 0; j < count && i-j > 0; j++ {
							if isit(result[i-j-1]) {
								count++

							}

							result[i-j-1] = capitalize(result[i-j-1])
						}

						result[i] = ""
						continue

					}
				}
			}
		}
	}

	filteredResult := make([]string, 0)
	for _, word := range result {
		if word != "" {
			filteredResult = append(filteredResult, word)
		}
	}

	return strings.Join(filteredResult, " ")
}

func FormatFlags(input string) string {
	re := regexp.MustCompile(`\(\w+,\s*-?\d+\)`)
	formatted := re.ReplaceAllStringFunc(input, func(match string) string {
		return regexp.MustCompile(`\s+`).ReplaceAllString(match, "")
	})
	return formatted
}

func isit(r string) bool {
	return r == "." || r == "," || r == "!" || r == "?" || r == ":" || r == ";" || r == "'" || r == ""

}
