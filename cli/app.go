package cli

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/eggsampler/certgot/log"
)

var (
	regArgLong  = regexp.MustCompile(`^--([[:alnum:]]+[\-[[:alnum:]]]*)(?:=(.+))?$`)
	regArgShort = regexp.MustCompile("^-(a+|b+|c+|d+|e+|f+|g+|h+|i+|j+|k+|l+|m+|n+|o+|p+|q+|r+|s+|t+|u+|v+|w+|x+|y+|z+)(?:=(.+))?$")
)

func extractArg(s string) []string {
	if strings.HasPrefix(s, "--") {
		return regArgLong.FindStringSubmatch(s)
	}
	return regArgShort.FindStringSubmatch(s)
}

type App struct {
	Name string

	FuncPreRun      func(*App) error
	FuncPostRun     func(*App, interface{})
	FuncRecover     func(*App, interface{})
	FuncHelpPrinter func(*App)

	argsMap  map[string]*Argument
	argsList []*Argument

	subCommandMap     map[string]*SubCommand
	subCommandList    []*SubCommand
	defaultSubCommand string

	helpTopics []HelpTopic

	// context
	OriginalArgs       []string
	SpecificSubCommand *SubCommand

	// mainly for testing
	exitFunc func(int)
}

func (app App) Exit(i int) {
	if app.exitFunc != nil {
		app.exitFunc(i)
		return
	}

	os.Exit(i)
}

func (app *App) GetArguments() map[string]*Argument {
	return app.argsMap
}

// TODO: not sure this is super helpful
// ie, should apps specifically store their args as variables
// and refer to them by the variable
// rather than this getter func
func (app *App) GetArgument(key string) *Argument {
	return app.argsMap[key]
}

func (app *App) GetSubCommands() map[string]*SubCommand {
	return app.subCommandMap
}

func (app *App) AddArgument(argument *Argument) {
	// TODO: error/panic if argument name contains space?
	if argument == nil {
		return
	}
	app.argsList = append(app.argsList, argument)
	if app.argsMap == nil {
		app.argsMap = map[string]*Argument{}
	}
	app.argsMap[argument.Name] = argument
	for _, altName := range argument.AltNames {
		app.argsMap[altName] = argument
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
	app.subCommandList = append(app.subCommandList, subCommand)
	if app.subCommandMap == nil {
		app.subCommandMap = map[string]*SubCommand{}
	}
	app.subCommandMap[subCommand.Name] = subCommand
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

func (app *App) Run() error {
	calledPostRun := false
	var postRunBlah interface{}

	if app.FuncPostRun != nil {
		defer func() {
			if !calledPostRun {
				calledPostRun = true
				app.FuncPostRun(app, postRunBlah)
			}
		}()
	}

	if app.FuncRecover != nil {
		// only recover if recover func set, otherwise lets panic fall through
		defer func() {
			if r := recover(); r != nil {
				postRunBlah = r
				app.FuncRecover(app, r)
			}
		}()
	}

	if app.FuncPreRun != nil {
		if err := app.FuncPreRun(app); err != nil {
			postRunBlah = err
			return err
		}
	}

	if app.SpecificSubCommand != nil && app.SpecificSubCommand.Run != nil {
		if err := app.SpecificSubCommand.Run(app); err != nil {
			postRunBlah = err
			return err
		}
	} else {
		app.PrintHelp()
	}

	return nil
}

func (app *App) Parse(argsToParse []string) error {
	// TODO
	/*
		if len(argsToParse) == 0 {
			app.PrintHelp()
			return nil
		}
	*/

	for _, arg := range app.argsMap {
		if arg.PreParse != nil {
			if err := arg.PreParse(arg, app); err != nil {
				return fmt.Errorf("error in argument %q PreParse func: %w", arg.Name, err)
			}
		}
	}

	argLast := ""

	for argIdx, v := range argsToParse {
		if strings.HasPrefix(v, "-") {
			argMatch := extractArg(v)
			if len(argMatch) < 2 {
				return fmt.Errorf("invalid argument: %s", argsToParse[argIdx])
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

			arg := app.argsMap[argLast]
			if arg == nil {
				return fmt.Errorf("unknown argument: %s", v)
			}

			if !arg.IsPresent && arg.OnPresent != nil {
				if err := arg.OnPresent(arg, argLast, argCount, app); err != nil {
					return fmt.Errorf("error in argument %q OnPresent func: %w", arg.Name, err)
				}
			}
			arg.IsPresent = true
			arg.RepeatCount = argCount

			if strings.Contains(argMatch[0], "=") {
				argLast = ""

				if arg.OnSet != nil {
					if err := arg.OnSet(arg, argLast, argMatch[2], app); err != nil {
						return fmt.Errorf("error in argument %q OnSet func: %w", arg.Name, err)
					}
				}

				if err := arg.Set(argMatch[2]); err != nil {
					return fmt.Errorf("error setting inline argument %s: %w", v, err)
				}
			}
		} else if argLast != "" {
			arg := app.argsMap[argLast]
			if arg.OnSet != nil {
				if err := arg.OnSet(arg, argLast, v, app); err != nil {
					return fmt.Errorf("error in argument %q OnSet func: %w", arg.Name, err)
				}
			}
			if err := arg.Set(v); err != nil {
				return fmt.Errorf("error setting argument %q: %w", argLast, err)
			}
			argLast = ""
		} else if app.SpecificSubCommand != nil {
			return fmt.Errorf("extra subcommand %q found, already provided %q", v, app.SpecificSubCommand.Name)
		} else {
			app.SpecificSubCommand = app.subCommandMap[v]
			if app.SpecificSubCommand == nil {
				return fmt.Errorf("invalid subcommand: %s", v)
			}
		}
	}

	if app.SpecificSubCommand == nil {
		app.SpecificSubCommand = app.subCommandMap[app.defaultSubCommand]
	}

	for _, arg := range app.argsMap {
		if arg.PostParse != nil {
			if err := arg.PostParse(arg, app.SpecificSubCommand, app); err != nil {
				if errors.Is(err, ErrExitSuccess) {
					log.WithError(err).WithField("argName", arg.Name).Debug("PostParse returned ErrExitSuccess")
					app.Exit(0)
					return nil
				}
				return fmt.Errorf("error in argument %q PostParse func: %w", arg.Name, err)
			}
		}
	}

	return nil
}

func (app *App) LoadConfig(cfgFile *Argument) error {
	return loadConfig(app, cfgFile, os.DirFS(""))
}

func (app *App) PrintHelp(topic ...string) {
	// TODO: not sure if this is useful, topic isn't used yet, will it be used??
	if len(topic) > 0 {
		helpArg := app.GetArgument("help")
		if helpArg != nil {
			_ = helpArg.Set(topic[0])
		}
	}

	if app.FuncHelpPrinter != nil {
		app.FuncHelpPrinter(app)
		return
	}

	DefaultHelpPrinter(app)
}

func (app *App) AddHelpTopic(topic HelpTopic) {
	app.helpTopics = append(app.helpTopics, topic)
}

func (app *App) AddHelpTopics(topics ...HelpTopic) {
	for _, v := range topics {
		app.AddHelpTopic(v)
	}
}

func (app *App) GetHelpTopics() []HelpTopic {
	return app.helpTopics
}
