package cli

import (
	"os"

	"golang.org/x/term"
)

var (
	termWidth  = 80
	isTerminal = false
)

func init() {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		termWidth = w
	}
	isTerminal = term.IsTerminal(int(os.Stdout.Fd()))
}

func TermWidth() int {
	return termWidth
}

func IsTerminal() bool {
	return isTerminal
}

func contains(s []string, c string) bool {
	for _, v := range s {
		if v == c {
			return true
		}
	}
	return false
}
