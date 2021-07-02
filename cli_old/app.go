package cli_old

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/eggsampler/certgot/log"
)

var (
	regFlagLong  = regexp.MustCompile(`^--([[:alnum:]]+(?:-[[:alnum:]]+)*)(?:=(.+))?$`)
	regFlagShort = regexp.MustCompile("^-(a+|b+|c+|d+|e+|f+|g+|h+|i+|j+|k+|l+|m+|n+|o+|p+|q+|r+|s+|t+|u+|v+|w+|x+|y+|z+)(?:=(.+))?$")
)

func extractFlag(s string) []string {
	if strings.HasPrefix(s, "--") {
		return regFlagLong.FindStringSubmatch(s)
	}
	return regFlagShort.FindStringSubmatch(s)
}

type App struct {
	Name string

	FuncPreRun      func(*App) error
	FuncPostRun     func(*App, interface{})
	FuncRecover     func(*App, interface{})
	FuncHelpPrinter func(*App, string)

	flagsMap  map[string]*Flag
	flagsList []*Flag

	subCommandMap     map[string]*SubCommand
	subCommandList    []*SubCommand
	defaultSubCommand string

	helpTopics []HelpTopic

	// context
	OriginalArguments []string
	FoundSubCommand   *SubCommand
	ExtraArguments    []string

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

func (app *App) Flags() map[string]*Flag {
	return app.flagsMap
}

// TODO: not sure this is super helpful
// ie, should apps specifically store their args as variables
// and refer to them by the variable
// rather than this getter func
func (app *App) Flag(key string) *Flag {
	return app.flagsMap[key]
}

func (app *App) SubCommands() map[string]*SubCommand {
	return app.subCommandMap
}

func (app *App) SubCommand(key string) *SubCommand {
	return app.subCommandMap[key]
}

func (app *App) AddFlag(f *Flag) {
	// TODO: error/panic if f name contains space?
	if f == nil {
		return
	}
	app.flagsList = append(app.flagsList, f)
	if app.flagsMap == nil {
		app.flagsMap = map[string]*Flag{}
	}
	app.flagsMap[f.Name] = f
	for _, altName := range f.AltNames {
		app.flagsMap[altName] = f
	}
}

func (app *App) AddFlags(f ...*Flag) {
	for _, v := range f {
		app.AddFlag(v)
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

	if app.FoundSubCommand != nil && app.FoundSubCommand.Run != nil {
		if err := app.FoundSubCommand.Run(app); err != nil {
			postRunBlah = err
			return err
		}
	} else {
		app.PrintHelp("")
	}

	return nil
}

func flagDashes(s string) string {
	if len(s) == 1 {
		return "-"
	}
	return "--"
}

func joinFlagNames(f Flag) string {
	s := []string{
		flagDashes(f.Name) + f.Name,
	}
	for _, n := range f.AltNames {
		s = append(s, flagDashes(n)+n)
	}
	return strings.Join(s, "/")
}

func (app *App) Parse(argsToParse []string) error {
	// TODO
	/*
		if len(argsToParse) == 0 {
			app.PrintHelp()
			return nil
		}
	*/

	for _, f := range app.flagsMap {
		if f.PreParse != nil {
			if err := f.PreParse(f, app); err != nil {
				return fmt.Errorf("error in argument %q PreParse func: %w", joinFlagNames(*f), err)
			}
		}
	}

	argLast := ""

	for argIdx, v := range argsToParse {
		if strings.HasPrefix(v, "-") {
			flagMatch := extractFlag(v)
			if len(flagMatch) < 2 {
				return fmt.Errorf("invalid argument: %s", argsToParse[argIdx])
			}

			argLast = flagMatch[1]
			var argCount int
			if !strings.HasPrefix(v, "--") {
				argCount = len(argLast)
				argLast = string(argLast[0])
			} else {
				argCount = 1
				argLast = flagMatch[1]
			}

			f := app.flagsMap[argLast]
			if f == nil {
				return fmt.Errorf("unknown argument: %s", v)
			}

			if !f.isPresent && f.OnPresent != nil {
				if err := f.OnPresent(f, argLast, argCount, app); err != nil {
					return fmt.Errorf("error in argument %q OnPresent func: %w", joinFlagNames(*f), err)
				}
			}
			f.isPresent = true
			f.isPresentInArgument = true
			f.RepeatCount = argCount

			if strings.Contains(flagMatch[0], "=") {
				argLast = ""

				if f.OnSet != nil {
					if err := f.OnSet(f, argLast, flagMatch[2], app); err != nil {
						return fmt.Errorf("error in argument %q OnSet func: %w", joinFlagNames(*f), err)
					}
				}

				if err := f.Set(flagMatch[2]); err != nil {
					return fmt.Errorf("error setting inline argument %s: %w", v, err)
				}
			}
		} else if argLast != "" {
			f := app.flagsMap[argLast]
			if f.OnSet != nil {
				if err := f.OnSet(f, argLast, v, app); err != nil {
					return fmt.Errorf("error in argument %q OnSet func: %w", joinFlagNames(*f), err)
				}
			}
			if err := f.Set(v); err != nil {
				return fmt.Errorf("error setting argument %q: %w", argLast, err)
			}
			argLast = ""
		} else if app.FoundSubCommand != nil {
			// anything after a subcommand has been found is added to the extra arguments
			app.ExtraArguments = append(app.ExtraArguments, v)
		} else {
			// if no subcommand has been found, attempt to locate it
			app.FoundSubCommand = app.subCommandMap[v]
			if app.FoundSubCommand == nil {
				return fmt.Errorf("invalid subcommand: %s", v)
			}
		}
	}

	if app.FoundSubCommand == nil {
		app.FoundSubCommand = app.subCommandMap[app.defaultSubCommand]
	}

	for _, f := range app.flagsMap {
		if f.PostParse != nil {
			if err := f.PostParse(f, app.FoundSubCommand, app); err != nil {
				if errors.Is(err, ErrExitSuccess) {
					log.WithError(err).WithField("flagName", f.Name).Debug("PostParse returned ErrExitSuccess")
					app.Exit(0)
					return nil
				}
				return fmt.Errorf("error in argument %q PostParse func: %w", joinFlagNames(*f), err)
			}
		}
	}

	return nil
}

func (app *App) LoadConfig(cfgFile *Flag) error {
	return loadConfig(app, cfgFile, os.DirFS(""))
}

func (app *App) PrintHelp(topic string) {
	if app.FuncHelpPrinter != nil {
		app.FuncHelpPrinter(app, topic)
		return
	}

	DefaultHelpPrinter(app, topic)
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
