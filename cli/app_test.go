package cli

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func Test_extractArg(t *testing.T) {
	testList := []struct {
		testName string
		arg      string
		matches  bool
	}{
		{
			testName: "not an arg",
			arg:      "asdf",
			matches:  false,
		},
		{
			testName: "too many dashes",
			arg:      "---asdf",
			matches:  false,
		},
		{
			testName: "long arg form",
			arg:      "--asdf",
			matches:  true,
		},
		{
			testName: "single short arg form",
			arg:      "-a",
			matches:  true,
		},
		{
			testName: "multiple short arg form",
			arg:      "-aaa",
			matches:  true,
		},
		{
			testName: "multiple short arg form with value",
			arg:      "-aaa=asd",
			matches:  true,
		},
		{
			testName: "multiple diff short arg form",
			arg:      "-asdf",
			matches:  false,
		},
	}

	for _, currentTest := range testList {
		t.Run(currentTest.testName, func(t *testing.T) {
			m := extractArg(currentTest.arg)
			if currentTest.matches && m == nil {
				t.Fatalf("test %q: expected match for arg %q, got none", currentTest.testName, currentTest.arg)
			}
			if !currentTest.matches && m != nil {
				t.Fatalf("test %q: expected no match for arg %q, got: %s", currentTest.testName, currentTest.arg, m)
			}
		})
	}
}

func TestApp_AddArguments(t *testing.T) {
	testArg := &Argument{Name: "test1"}
	testArgShort := &Argument{Name: "test2", AltNames: []string{"t"}}

	testList := []struct {
		testName      string
		arguments     []*Argument
		expectedCount int
	}{
		{
			testName: "no arguments",
		},
		{
			testName:  "empty arguments",
			arguments: []*Argument{},
		},
		{
			testName:  "nil argument",
			arguments: []*Argument{nil},
		},
		{
			testName:      "single simple argument",
			arguments:     []*Argument{testArg},
			expectedCount: 1,
		},
		{
			testName:      "single less simple argument",
			arguments:     []*Argument{testArgShort},
			expectedCount: 2,
		},
		{
			testName:      "multiple arguments",
			arguments:     []*Argument{testArg, testArgShort},
			expectedCount: 3,
		},
	}

	for _, currentTest := range testList {
		t.Run(currentTest.testName, func(t *testing.T) {
			cli := App{}
			cli.AddArguments(currentTest.arguments...)
			args := cli.GetArguments()
			if len(args) != currentTest.expectedCount {
				t.Fatalf("test %q: expected %d arguments, got: %d",
					currentTest.testName, currentTest.expectedCount, len(cli.argsMap))
			}

			keyMap := map[string]bool{}
			for k := range args {
				keyMap[k] = false
			}
			for _, v := range currentTest.arguments {
				if v == nil {
					continue
				}
				_, ok := args[v.Name]
				if !ok {
					t.Fatalf("test %q: args not equal, doesn't contain arg name: %q", currentTest.testName, v.Name)
				}
				keyMap[v.Name] = true
				for _, vv := range v.AltNames {
					_, ok = args[vv]
					if !ok {
						t.Fatalf("test %q: args not equal, doesn't contain arg alt name: %q",
							currentTest.testName, v.Name)
					}
					keyMap[vv] = true
				}
			}
			for k, v := range keyMap {
				if !v {
					t.Fatalf("test %q: args not equal, didn't contain arg name: %q", currentTest.testName, k)
				}
			}
		})
	}
}

func TestApp_AddSubCommands(t *testing.T) {
	testList := []struct {
		testName      string
		subCommands   []*SubCommand
		expectedCount int
	}{
		{
			testName: "no subcommands",
		},
		{
			testName:    "empty subcommands",
			subCommands: []*SubCommand{},
		},
		{
			testName:    "nil subcommands",
			subCommands: []*SubCommand{nil},
		},
	}

	for _, currentTest := range testList {
		t.Run(currentTest.testName, func(t *testing.T) {
			cli := App{}
			cli.AddSubCommands(currentTest.subCommands...)
			subCmds := cli.GetSubCommands()
			if len(subCmds) != currentTest.expectedCount {
				t.Fatalf("test %q: expected %d arguments, got: %d",
					currentTest.testName, currentTest.expectedCount, len(cli.argsMap))
			}

			keyMap := map[string]bool{}
			for k := range subCmds {
				keyMap[k] = false
			}
			for _, v := range currentTest.subCommands {
				if v == nil {
					continue
				}
				_, ok := subCmds[v.Name]
				if !ok {
					t.Fatalf("test %q: subcommands not equal, doesn't contain subcommands name: %q",
						currentTest.testName, v.Name)
				}
				keyMap[v.Name] = true
			}
			for k, v := range keyMap {
				if !v {
					t.Fatalf("test %q: subcommands not equal, didn't contain subcommands name: %q",
						currentTest.testName, k)
				}
			}
		})
	}
}

func TestApp_Run(t *testing.T) {
	testList := []struct {
		testName   string
		app        *App
		subCommand *SubCommand
		hasError   bool
		errorStr   string
	}{
		{
			testName: "empty app",
			app:      &App{},
		},
		{
			testName: "app recover",
			app: &App{
				FuncRecover: func(app *App, r interface{}) {
					s := fmt.Sprintf("%s", r)
					if s != "HELLO WORLD!" {
						panic("unknown panic recovered")
					}
				},
			},
			subCommand: &SubCommand{Run: func(app *App) error {
				panic("HELLO WORLD!")
			}},
		},
		{
			testName: "app recover postrun 1",
			app: &App{
				FuncPostRun: func(app *App, r interface{}) {
					s := fmt.Sprintf("%s", r)
					if s != "HELLO WORLD!" {
						panic("unknown panic message: " + s)
					}
				},
				FuncRecover: func(app *App, r interface{}) {
					s := fmt.Sprintf("%s", r)
					if s != "HELLO WORLD!" {
						panic("unknown panic recovered: " + s)
					}
				},
			},
			subCommand: &SubCommand{Run: func(app *App) error {
				panic("HELLO WORLD!")
			}},
		},
		{
			testName: "app recover postrun 2",
			app: &App{
				FuncPostRun: func(app *App, r interface{}) {
					s := fmt.Sprintf("%s", r)
					if s != "HELLO WORLD!" {
						panic("unknown panic message: " + s)
					}
				},
				FuncRecover: func(app *App, r interface{}) {
					s := fmt.Sprintf("%s", r)
					if s != "HELLO WORLD!" {
						panic("unknown panic recovered: " + s)
					}
				},
			},
			subCommand: &SubCommand{Run: func(app *App) error {
				panic("HELLO WORLD!")
			}},
		},
		{
			testName: "prefunc succeed",
			app: &App{
				FuncPreRun: func(app *App) error {
					return nil
				},
			},
		},
		{
			testName: "prefunc fail",
			app: &App{
				FuncPreRun: func(app *App) error {
					return errors.New("prefunc funky")
				},
			},
			hasError: true,
			errorStr: "prefunc funky",
		},
		{
			testName: "SubCommand succeed",
			app:      &App{},
			subCommand: &SubCommand{
				Run: func(app *App) error {
					return nil
				},
			},
		},
		{
			testName: "SubCommand fail",
			app:      &App{},
			subCommand: &SubCommand{
				Run: func(app *App) error {
					return errors.New("subcommand funky")
				},
			},
			hasError: true,
			errorStr: "subcommand funky",
		},
		{
			testName: "postfunc",
			app: &App{
				FuncPostRun: func(app *App, r interface{}) {
					// hello world
				},
			},
		},
	}

	for _, currentTest := range testList {
		t.Run(currentTest.testName, func(t *testing.T) {
			if currentTest.subCommand != nil {
				currentTest.app.AddSubCommand(currentTest.subCommand)
				currentTest.app.SpecificSubCommand = currentTest.subCommand
			}
			err := currentTest.app.Run()
			if currentTest.hasError == (err == nil) {
				t.Fatalf("%q: expected error %v, got: %v", currentTest.testName, currentTest.hasError, err)
			}
			if err != nil && !strings.Contains(err.Error(), currentTest.errorStr) {
				t.Fatalf("test %q: expected %q in error: %v", currentTest.testName, currentTest.errorStr, err)
			}
		})
	}
}

func Test_doParse(t *testing.T) {
	testList := []struct {
		testName       string
		appSubCommands []*SubCommand
		appArguments   []*Argument
		argsToParse    []string
		hasError       bool
		errorStr       string
		checkFunc      func(*SubCommand, map[string]*Argument) error
	}{
		{
			testName:    "malformed argument",
			argsToParse: []string{"-"},
			hasError:    true,
			errorStr:    "invalid argument",
		},
		{
			testName:    "unknown argument no argument",
			argsToParse: []string{"--unknown"},
			hasError:    true,
			errorStr:    "unknown argument",
		},
		{
			testName:     "unknown argument partial",
			appArguments: []*Argument{{Name: "unknown"}},
			argsToParse:  []string{"--unknown1"},
			hasError:     true,
			errorStr:     "unknown argument",
		},
		{
			testName:       "default subcommand",
			appSubCommands: []*SubCommand{{Name: "default", Default: true}},
			checkFunc: func(sc *SubCommand, args map[string]*Argument) error {
				if sc.Name != "default" {
					return fmt.Errorf("expected default, got: %s", sc.Name)
				}
				return nil
			},
		},
		{
			testName: "non-default subcommand",
			appSubCommands: []*SubCommand{
				{Name: "default", Default: true},
				{Name: "non-default", Default: true},
			},
			argsToParse: []string{"non-default"},
			checkFunc: func(sc *SubCommand, args map[string]*Argument) error {
				if sc.Name != "non-default" {
					return fmt.Errorf("expected non-default, got: %s", sc.Name)
				}
				return nil
			},
		},
		{
			testName:       "known arg not present",
			appSubCommands: []*SubCommand{{Name: "default", Default: true}},
			appArguments:   []*Argument{{Name: "known"}},
			argsToParse:    []string{"--known"},
			checkFunc: func(sc *SubCommand, args map[string]*Argument) error {
				knownArg := args["known"]
				if !knownArg.isPresent {
					return fmt.Errorf("expected known arg is not present")
				}
				return nil
			},
		},
		{
			testName:       "known arg = no value",
			appSubCommands: []*SubCommand{{Name: "default", Default: true}},
			appArguments:   []*Argument{{Name: "known"}},
			argsToParse:    []string{"--known=value"},
			hasError:       true,
			errorStr:       "does not take a value",
		},
		{
			testName:     "known arg no value",
			appArguments: []*Argument{{Name: "known"}},
			argsToParse:  []string{"--known", "value"},
			hasError:     true,
			errorStr:     "does not take a value",
		},
		{
			testName:       "known arg = value",
			appSubCommands: []*SubCommand{{Name: "default", Default: true}},
			appArguments:   []*Argument{{Name: "known", TakesValue: true}},
			argsToParse:    []string{"--known=value"},
			checkFunc: func(sc *SubCommand, args map[string]*Argument) error {
				knownArg := args["known"]
				if knownArg.String() != "value" {
					return fmt.Errorf("expected \"value\", got: %q", knownArg.String())
				}
				return nil
			},
		},
		{
			testName:       "known arg value",
			appSubCommands: []*SubCommand{{Name: "default", Default: true}},
			appArguments:   []*Argument{{Name: "known", TakesValue: true}},
			argsToParse:    []string{"--known", "value"},
			checkFunc: func(sc *SubCommand, args map[string]*Argument) error {
				knownArg := args["known"]
				if knownArg.String() != "value" {
					return fmt.Errorf("expected \"value\", got: %q", knownArg.String())
				}
				return nil
			},
		},
		{
			testName:       "extra same subcommand",
			appSubCommands: []*SubCommand{{Name: "default", Default: true}},
			argsToParse:    []string{"default", "default"},
			hasError:       true,
			errorStr:       "extra subcommand",
		},
		{
			testName: "extra different subcommand",
			appSubCommands: []*SubCommand{
				{Name: "default", Default: true},
				{Name: "extra"},
			},
			argsToParse: []string{"default", "extra"},
			hasError:    true,
			errorStr:    "extra subcommand",
		},
		{
			testName: "non-default subcommand with argument and value",
			appSubCommands: []*SubCommand{
				{Name: "default", Default: true},
				{Name: "extra"},
			},
			appArguments: []*Argument{{Name: "known", TakesValue: true}},
			argsToParse:  []string{"--known", "value", "extra"},
			checkFunc: func(sc *SubCommand, args map[string]*Argument) error {
				if sc.Name != "extra" {
					return fmt.Errorf("expected extra, got: %v", sc.Name)
				}
				return nil
			},
		},
		{
			testName: "preparse func fail",
			appArguments: []*Argument{{Name: "hello", PreParse: func(arg *Argument, app *App) error {
				return errors.New("hello world")
			}}},
			hasError: true,
			errorStr: "hello world",
		},
		{
			testName: "preparse func succeed",
			appArguments: []*Argument{{Name: "hello", PreParse: func(arg *Argument, app *App) error {
				return nil
			}}},
		},
		{
			testName: "onpresent func fail",
			appArguments: []*Argument{
				{
					Name: "hello",
					OnPresent: func(arg *Argument, argString string, repeatCount int, app *App) error {
						return errors.New("hello world")
					},
				},
			},
			argsToParse: []string{"--hello"},
			hasError:    true,
			errorStr:    "hello world",
		},
		{
			testName: "onpresent func succeed",
			appArguments: []*Argument{{Name: "hello", PreParse: func(arg *Argument, app *App) error {
				return nil
			}}},
			argsToParse: []string{"--hello"},
		},
		{
			testName: "onpresent func not present",
			appArguments: []*Argument{
				{Name: "notused"},
				{Name: "hello", OnPresent: func(arg *Argument, argString string, repeatCount int, app *App) error {
					return errors.New("hello world")
				}}},
			argsToParse: []string{"--notused"},
		},
		{
			testName:     "args name repeated",
			appArguments: []*Argument{{Name: "v"}},
			argsToParse:  []string{"-vvvvv"},
			checkFunc: func(sc *SubCommand, args map[string]*Argument) error {
				arg := args["v"]
				if arg.RepeatCount != 5 {
					return fmt.Errorf("arg count mismatched, expected 5, got: %d", arg.RepeatCount)
				}
				return nil
			},
		},
		{
			testName: "onset 1",
			appArguments: []*Argument{{Name: "testonset", OnSet: func(arg *Argument, argString string, newValue interface{}, app *App) error {
				if newValue == "hello" {
					return errors.New("hello set!")
				}
				return fmt.Errorf("unknown set: %v", newValue)
			}}},
			argsToParse: []string{"--testonset=hello"},
			hasError:    true,
			errorStr:    "hello set!",
		},
		{
			testName: "onset 2",
			appArguments: []*Argument{{Name: "testonset", OnSet: func(arg *Argument, argString string, newValue interface{}, app *App) error {
				if newValue == "hello" {
					return errors.New("hello set!")
				}
				return fmt.Errorf("unknown set: %v", newValue)
			}}},
			argsToParse: []string{"--testonset", "hello"},
			hasError:    true,
			errorStr:    "hello set!",
		},
		{
			testName:    "invalid subcommand",
			argsToParse: []string{"nonexistant"},
			hasError:    true,
			errorStr:    "invalid subcommand",
		},
		{
			testName: "postparse func fail",
			appArguments: []*Argument{{Name: "hello", PostParse: func(arg *Argument, sc *SubCommand, app *App) error {
				return errors.New("hello world")
			}}},
			hasError: true,
			errorStr: "hello world",
		},
		{
			testName: "postparse func succeed",
			appArguments: []*Argument{{Name: "hello", PostParse: func(arg *Argument, sc *SubCommand, app *App) error {
				return nil
			}}},
		},
	}

	for _, currentTest := range testList {
		t.Run(currentTest.testName, func(t *testing.T) {
			app := App{}
			app.AddSubCommands(currentTest.appSubCommands...)
			app.AddArguments(currentTest.appArguments...)
			sc, err := doParse(&app, currentTest.argsToParse)
			if currentTest.hasError == (err == nil) {
				t.Fatalf("%q: expected error %v, got: %v", currentTest.testName, currentTest.hasError, err)
			}
			if err != nil && !strings.Contains(err.Error(), currentTest.errorStr) {
				t.Fatalf("test %q: expected %q in error: %v", currentTest.testName, currentTest.errorStr, err)
			}
			if currentTest.checkFunc != nil {
				if err := currentTest.checkFunc(sc, app.argsMap); err != nil {
					t.Fatalf("test %q: check: %v", currentTest.testName, err)
				}
			}
		})
	}
}
