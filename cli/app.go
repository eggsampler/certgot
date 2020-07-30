package cli

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	regArgLong  = regexp.MustCompile(`^--([[:alnum:]]+[\-[[:alnum:]]]*)(?:=(.+))?$`)
	regArgShort = regexp.MustCompile("^-(a+|b+|c+|d+|e+|f+|g+|h+|i+|j+|k+|l+|m+|n+|o+|p+|q+|r+|s+|t+|u+|v+|w+|x+|y+|z+)(?:=(.+))?$")
)

type App struct {
	FuncPreRun      func(c *Context) error
	FuncPostRun     func(c *Context) error
	FuncRecover     func(c *Context, r interface{}) error
	FuncHelpPrinter func(app *App)

	args              map[string]*Argument
	subCommands       map[string]*SubCommand
	defaultSubCommand string
}

func (app *App) GetArguments() map[string]*Argument {
	return app.args
}

func (app *App) GetArgument(key string) *Argument {
	return app.args[key]
}

func (app *App) GetSubCommands() map[string]*SubCommand {
	return app.subCommands
}

func (app *App) AddArgument(argument *Argument) {
	// TODO: error/panic if argument name contains space?
	if argument == nil {
		return
	}
	if app.args == nil {
		app.args = map[string]*Argument{}
	}
	app.args[argument.Name] = argument
	for _, altName := range argument.AltNames {
		app.args[altName] = argument
	}
}

func (app *App) AddArguments(arguments ...*Argument) {
	for _, v := range arguments {
		app.AddArgument(v)
	}
}

func (app *App) AddSubCommand(subCommand *SubCommand) {
	if subCommand == nil {
		return
	}
	if app.subCommands == nil {
		app.subCommands = map[string]*SubCommand{}
	}
	app.subCommands[subCommand.Name] = subCommand
	if subCommand.Default {
		app.defaultSubCommand = subCommand.Name
	}
}

func (app *App) AddSubCommands(commands ...*SubCommand) {
	if commands == nil {
		return
	}
	for _, v := range commands {
		app.AddSubCommand(v)
	}
}

func doRun(ctx *Context) error {
	calledPostRun := false

	if ctx.App.FuncRecover != nil {
		// only recover if recover func set, otherwise lets panic fall through
		defer func() {
			if r := recover(); r != nil {
				_ = ctx.App.FuncRecover(ctx, r)
			}
			if !calledPostRun && ctx.App.FuncPostRun != nil {
				calledPostRun = true
				_ = ctx.App.FuncPostRun(ctx)
			}
		}()
	}

	if ctx.App.FuncPreRun != nil {
		if err := ctx.App.FuncPreRun(ctx); err != nil {
			return err
		}
	}

	if ctx.SubCommand.Run != nil {
		if err := ctx.SubCommand.Run(ctx); err != nil {
			return err
		}
	}

	if ctx.App.FuncPostRun != nil {
		calledPostRun = true
		if err := ctx.App.FuncPostRun(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) Run(argsToParse []string) error {
	if len(argsToParse) == 0 {
		app.PrintHelp()
		return nil
	}

	sc, err := app.parse(argsToParse)
	if err != nil {
		return err
	}

	ctx := Context{
		OriginalArguments: argsToParse,
		App:               app,
		SubCommand:        *sc,
	}

	return doRun(&ctx)
}

func extractArg(s string) []string {
	if strings.HasPrefix(s, "--") {
		return regArgLong.FindStringSubmatch(s)
	}
	return regArgShort.FindStringSubmatch(s)
}

func (app *App) parse(argsToParse []string) (*SubCommand, error) {
	for _, arg := range app.args {
		if arg.PreParse != nil {
			if err := arg.PreParse(arg, app); err != nil {
				return nil, fmt.Errorf("error in argument %q PreParse func: %w", arg.Name, err)
			}
		}
	}

	var sc *SubCommand
	argLast := ""

	for argIdx, v := range argsToParse {
		if strings.HasPrefix(v, "-") {
			argMatch := extractArg(v)
			if len(argMatch) < 2 {
				return sc, fmt.Errorf("invalid argument: %s", argsToParse[argIdx])
			}

			argLast = argMatch[1]
			var argCount int
			if !strings.HasPrefix(v, "--") {
				argCount = len(argLast)
				argLast = string(argLast[0])
			} else {
				argCount = 1
				argLast = argMatch[1]
			}

			arg := app.args[argLast]
			if arg == nil {
				return sc, fmt.Errorf("unknown argument: %s", v)
			}

			if !arg.isPresent && arg.OnPresent != nil {
				if err := arg.OnPresent(arg, argLast, argCount, app); err != nil {
					return sc, fmt.Errorf("error in argument %q OnPresent func: %w", arg.Name, err)
				}
			}
			arg.isPresent = true
			arg.RepeatCount = argCount

			if strings.Contains(argMatch[0], "=") {
				argLast = ""

				if arg.OnSet != nil {
					if err := arg.OnSet(arg, argLast, argMatch[2], app); err != nil {
						return sc, fmt.Errorf("error in argument %q OnSet func: %w", arg.Name, err)
					}
				}

				if err := arg.Set(argMatch[2]); err != nil {
					return sc, fmt.Errorf("error setting inline argument %s: %w", v, err)
				}
			}
		} else if argLast != "" {
			arg := app.args[argLast]
			if arg.OnSet != nil {
				if err := arg.OnSet(arg, argLast, v, app); err != nil {
					return sc, fmt.Errorf("error in argument %q OnSet func: %w", arg.Name, err)
				}
			}
			if err := arg.Set(v); err != nil {
				return sc, fmt.Errorf("error setting argument %s: %w", argLast, err)
			}
			argLast = ""
		} else if sc != nil {
			return sc, fmt.Errorf("extra subcommand %q found, already provided %q", v, sc.Name)
		} else {
			sc = app.subCommands[v]
			if sc == nil {
				return sc, fmt.Errorf("invalid subcommand: %s", v)
			}
		}
	}

	if sc == nil {
		sc = app.subCommands[app.defaultSubCommand]
	}

	for _, arg := range app.args {
		if arg.PostParse != nil {
			if err := arg.PostParse(arg, sc, app); err != nil {
				return sc, fmt.Errorf("error in argument %q PostParse func: %w", arg.Name, err)
			}
		}
	}

	return sc, nil
}

func (app *App) PrintHelp(topic ...string) {
	if app.FuncHelpPrinter != nil {
		app.FuncHelpPrinter(app)
		return
	}

	DefaultHelpPrinter(app)
}
