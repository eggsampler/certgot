package cli

import (
	"fmt"
	"os"
)

// App is the main entry point for any cli application, and should be defined as a variable
type App struct {
	// Name is the name of the app
	// Mostly just used for help purposes
	Name string

	// Flags is a list of flags which can be provided to the application
	Flags FlagList

	// Commands is a list of commands which can be run by the application
	Commands CommandList

	// Configs is a list of configurations to be set or used by the application
	Configs ConfigList

	// PreRunFunc is run before the Command.RunFunc
	// Can be used to do things like load config files, set up any common state, etc
	PreRunFunc func(*Context) error

	// PostRunFunc is run after the Command.RunFunc with the error, if any, returned
	PostRunFunc func(*Context, error) error

	// Help is a list of categories used when displaying help
	Help HelpCategories

	// HelpPrinterFunc can be used to override the default help printer
	// Could maybe be used to print help in json, or totally customise the help display, or even display no help at all
	HelpPrinterFunc func(ctx *Context, category string) error

	// RecoverFunc can be used to catch any panics during any point of the App.Run function
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

	// setup initial context
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
		if err := f.PostParseFunc(f, &ctx); err != nil {
			return fmt.Errorf("error on flag %s PostParseFunc: %w", f.Name, err)
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
		err = ctx.Command.RunFunc(&ctx)
	}

	// run the app post run
	if app.PostRunFunc != nil {
		if postErr := app.PostRunFunc(&ctx, err); err != nil {
			return fmt.Errorf("error in app PostRunFunc: %w", postErr)
		}
	}

	if err != nil {
		return fmt.Errorf("error in command %q RunFunc: %w", ctx.Command.Name, err)
	}

	// TODO: error if command is nil?
	// or a default app run func?

	return nil
}

func (app App) PrintHelp(ctx *Context, category string) error {
	if app.HelpPrinterFunc != nil {
		return app.HelpPrinterFunc(ctx, category)
	}

	DefaultHelpPrinter(ctx, category)

	return nil
}

func (app *App) LoadConfig(files []string, skip bool) error {
	return loadConfig(files, skip, app.Configs, os.DirFS(""))
}
