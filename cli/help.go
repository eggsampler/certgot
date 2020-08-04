package cli

import (
	"fmt"
	"strings"
)

func DefaultHelpPrinter(app *App) {
	fmt.Println(strings.Repeat("- ", 40))
	defer func() {
		fmt.Println(strings.Repeat("- ", 40))
	}()
	fmt.Println("TODO!")
}
