package log

import "strings"

func Wrap(s string, terminalLen int, prefix string) string {
	words := strings.Split(s, " ")

	var lines []string

	currentLine := ""
	for _, w := range words {
		if len(prefix)+len(currentLine)+1+len(w) >= terminalLen {
			lines = append(lines, prefix+currentLine)
			currentLine = w
			continue
		}
		if currentLine == "" {
			currentLine = w
			continue
		}
		currentLine += " " + w
	}
	if currentLine != "" {
		lines = append(lines, prefix+currentLine)
	}

	out := strings.Join(lines, "\n")
	return out
}
