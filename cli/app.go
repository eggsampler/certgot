package cli

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

	PreRunFunc  func(*Context) error
	PostRunFunc func(*Context) error

	Help            HelpCategories
	HelpPrinterFunc func(ctx *Context, category string) error

	RecoverFunc func(*App, interface{}) error
}

func (app *App) Run(args []string) (err error) {

	if app.RecoverFunc != nil {
		// only recover if recover func set, otherwise lets panic fall through
		defer func() {
			if r := recover(); r != nil {
				err = app.RecoverFunc(app, r)
			}
		}()
	}

	ctx := Context{
		App:          app,
		RawArguments: args,
	}

	// parse provided arguments
	if err := parseArguments(args, &ctx, app.Flags, app.Commands); err != nil {
		return fmt.Errorf("error parsing arguments: %w", err)
	}

	// execute any flag functions
	for _, f := range ctx.Flags {
		if err := f.OnSetFunc(f, &ctx); err != nil {
			return fmt.Errorf("error on flag %s OnSetFunc: %w", f.Name, err)
		}
	}

	// run the app pre run
	if app.PreRunFunc != nil {
		if err := app.PreRunFunc(&ctx); err != nil {
			return fmt.Errorf("error in app PreRunFunc: %w", err)
		}
	}

	// run the command
	if ctx.Command != nil {
		if err := ctx.Command.RunFunc(&ctx); err != nil {
			return fmt.Errorf("error in command %q RunFunc: %w", ctx.Command.Name, err)
		}
	}

	// run the app post run
	if app.PostRunFunc != nil {
		if err := app.PostRunFunc(&ctx); err != nil {
			return fmt.Errorf("error in app PostRunFunc: %w", err)
		}
	}

	// TODO: error if command is nil?
	// or a default app run func?

	return nil
}

func (app *App) PrintHelp(ctx *Context, category string) error {
	if app.HelpPrinterFunc != nil {
		return app.HelpPrinterFunc(ctx, category)
	}

	DefaultHelpPrinter(ctx, category)

	return nil
}
