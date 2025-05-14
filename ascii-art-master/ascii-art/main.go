package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . [string]")
		return
	}
	input := os.Args[1]

	input = strings.ReplaceAll(input, "\\n", "\n")

	bannerFile := "standard.txt"
	banner, err := readBanner(bannerFile)
	if err != nil {
		fmt.Println("Error reading banner:", err)
		return
	}
	lines := strings.Split(input, "\n") // تقسيم النص إلى أسطر بعد استبدال \\n بـ \n
	output := convertToASCIIWithDynamicSpaces(lines, banner)
	if isNewLine(input) {
		fmt.Print(output)
	} else {
		fmt.Println(output)

	}

}
func isNewLine(s string) bool {
	for _, chr := range s {
		if chr != '\n' {
			return false
		}
	}
	return true
}

func readBanner(filename string) (map[rune][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	banner := make(map[rune][]string)
	scanner := bufio.NewScanner(file)
	var currentChar rune
	charLines := []string{}
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if len(charLines) > 0 {
				banner[currentChar] = charLines
			}
			charLines = []string{}
			lineCount = 0
			continue
		}
		if lineCount == 0 {
			currentChar = rune(len(banner) + 32) // الحروف تبدأ من ASCII 32 (المسافة)
		}
		charLines = append(charLines, line)
		lineCount++
	}
	if len(charLines) > 0 {
		banner[currentChar] = charLines
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return banner, nil
}
func convertToASCIIWithDynamicSpaces(lines []string, banner map[rune][]string) string {
	var result []string
	emptyLineCount := 0
	for _, line := range lines {
		if line == "" {
			emptyLineCount++
			continue
		}
		if emptyLineCount > 0 {
			for i := 0; i < emptyLineCount; i++ {
				result = append(result, "")
			}
			emptyLineCount = 0
		}
		asciiLines := make([]string, 8)
		for _, char := range line {
			if charLines, exists := banner[char]; exists {
				for i := 0; i < 8; i++ {
					asciiLines[i] += charLines[i]
				}
			} else {

				fmt.Printf("ther is no %v in the stander.text ", string(char))
				return ""

			}
		}
		result = append(result, strings.Join(asciiLines, "\n"))

	}
	for i := 0; i < emptyLineCount; i++ {
		result = append(result, "")
	}
	return strings.Join(result, "\n")
}
