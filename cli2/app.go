package cli2

import (
	"fmt"
)

type App struct {
	// Name is the name of the app
	// Mostly just used for help purposes
	Name string

	Flags    FlagList
	Commands CommandList
	Configs  ConfigList

	Help        HelpCategories
	HelpPrinter func(app *App, category string) error
}

func (app *App) Run(args []string) error {

	ctx := Context{
		App:          app,
		RawArguments: args,
	}

	// parse provided arguments
	if err := parseArguments(args, &ctx, app.Flags, app.Commands); err != nil {
		return fmt.Errorf("error parsing arguments: %w", err)
	}

	// execute any flag functions

	return nil
}
