package cli_old

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_extractFlag(t *testing.T) {
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
			m := extractFlag(currentTest.arg)
			if currentTest.matches && m == nil {
				t.Fatalf("test %q: expected match for arg %q, got none", currentTest.testName, currentTest.arg)
			}
			if !currentTest.matches && m != nil {
				t.Fatalf("test %q: expected no match for arg %q, got: %s", currentTest.testName, currentTest.arg, m)
			}
		})
	}
}

func TestApp_AddFlag(t *testing.T) {
	testFlag := &Flag{Name: "test1"}
	testFlagShort := &Flag{Name: "test2", AltNames: []string{"t"}}

	testList := []struct {
		testName      string
		flags         []*Flag
		expectedCount int
	}{
		{
			testName: "no arguments",
		},
		{
			testName: "empty arguments",
			flags:    []*Flag{},
		},
		{
			testName: "nil argument",
			flags:    []*Flag{nil},
		},
		{
			testName:      "single simple argument",
			flags:         []*Flag{testFlag},
			expectedCount: 1,
		},
		{
			testName:      "single less simple argument",
			flags:         []*Flag{testFlagShort},
			expectedCount: 2,
		},
		{
			testName:      "multiple arguments",
			flags:         []*Flag{testFlag, testFlagShort},
			expectedCount: 3,
		},
	}

	for _, currentTest := range testList {
		t.Run(currentTest.testName, func(t *testing.T) {
			cli := App{}
			cli.AddFlags(currentTest.flags...)
			flags := cli.Flags()
			if len(flags) != currentTest.expectedCount {
				t.Fatalf("test %q: expected %d flags, got: %d",
					currentTest.testName, currentTest.expectedCount, len(cli.flagsMap))
			}

			keyMap := map[string]bool{}
			for k := range flags {
				keyMap[k] = false
			}
			for _, v := range currentTest.flags {
				if v == nil {
					continue
				}
				_, ok := flags[v.Name]
				if !ok {
					t.Fatalf("test %q: flags not equal, doesn't contain flag name: %q", currentTest.testName, v.Name)
				}
				keyMap[v.Name] = true
				for _, vv := range v.AltNames {
					_, ok = flags[vv]
					if !ok {
						t.Fatalf("test %q: flags not equal, doesn't contain flag alt name: %q",
							currentTest.testName, v.Name)
					}
					keyMap[vv] = true
				}
			}
			for k, v := range keyMap {
				if !v {
					t.Fatalf("test %q: flags not equal, didn't contain flag name: %q", currentTest.testName, k)
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
			subCmds := cli.SubCommands()
			if len(subCmds) != currentTest.expectedCount {
				t.Fatalf("test %q: expected %d subcommands, got: %d",
					currentTest.testName, currentTest.expectedCount, len(cli.flagsMap))
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
				currentTest.app.FoundSubCommand = currentTest.subCommand
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
		appFlags       []*Flag
		argsToParse    []string
		hasError       bool
		errorStr       string
		checkFunc      func(app *App, exitCalled bool, exitRet int) error
	}{
		{
			testName:    "malformed argument",
			argsToParse: []string{"-"},
			hasError:    true,
			errorStr:    "invalid argument",
		},
		{
			testName:    "unknown argument no argument 1",
			argsToParse: []string{"--unknown"},
			hasError:    true,
			errorStr:    "unknown argument",
		},
		{
			testName:    "unknown argument no argument 2",
			argsToParse: []string{"--un-known"},
			hasError:    true,
			errorStr:    "unknown argument",
		},
		{
			testName:    "unknown argument partial",
			appFlags:    []*Flag{{Name: "unknown"}},
			argsToParse: []string{"--unknown1"},
			hasError:    true,
			errorStr:    "unknown argument",
		},
		{
			testName:       "default subcommand",
			appSubCommands: []*SubCommand{{Name: "default", Default: true}},
			checkFunc: func(app *App, exitCalled bool, exitRet int) error {
				if app.FoundSubCommand.Name != "default" {
					return fmt.Errorf("expected default, got: %s", app.FoundSubCommand.Name)
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
			checkFunc: func(app *App, exitCalled bool, exitRet int) error {
				if app.FoundSubCommand.Name != "non-default" {
					return fmt.Errorf("expected non-default, got: %s", app.FoundSubCommand.Name)
				}
				return nil
			},
		},
		{
			testName:       "known arg not present",
			appSubCommands: []*SubCommand{{Name: "default", Default: true}},
			appFlags:       []*Flag{{Name: "known"}},
			argsToParse:    []string{"--known"},
			checkFunc: func(app *App, exitCalled bool, exitRet int) error {
				knownArg := app.flagsMap["known"]
				if !knownArg.isPresent {
					return fmt.Errorf("expected known arg is not present")
				}
				return nil
			},
		},
		{
			testName:       "known arg = no value",
			appSubCommands: []*SubCommand{{Name: "default", Default: true}},
			appFlags:       []*Flag{{Name: "known", Value: &SimpleValue{}}},
			argsToParse:    []string{"--known=value"},
			hasError:       true,
			errorStr:       "does not take a value",
		},
		{
			testName:    "known arg no value",
			appFlags:    []*Flag{{Name: "known", Value: &SimpleValue{}}},
			argsToParse: []string{"--known", "value"},
			hasError:    true,
			errorStr:    "does not take a value",
		},
		{
			testName:       "known arg = value",
			appSubCommands: []*SubCommand{{Name: "default", Default: true}},
			appFlags:       []*Flag{{Name: "known", TakesValue: true, Value: &SimpleValue{}}},
			argsToParse:    []string{"--known=value"},
			checkFunc: func(app *App, exitCalled bool, exitRet int) error {
				f := app.flagsMap["known"]
				s, err := f.String(false, false, false, false)
				if err != nil {
					return err
				}
				if s != "value" {
					return fmt.Errorf("expected \"value\", got: %q", s)
				}
				return nil
			},
		},
		{
			testName:       "known arg value",
			appSubCommands: []*SubCommand{{Name: "default", Default: true}},
			appFlags:       []*Flag{{Name: "known", TakesValue: true, Value: &SimpleValue{}}},
			argsToParse:    []string{"--known", "value"},
			checkFunc: func(app *App, exitCalled bool, exitRet int) error {
				f := app.flagsMap["known"]
				s, err := f.String(false, false, false, false)
				if err != nil {
					return err
				}
				if s != "value" {
					return fmt.Errorf("expected \"value\", got: %q", s)
				}
				return nil
			},
		},
		{
			testName:       "extra",
			appSubCommands: []*SubCommand{{Name: "default", Default: true}},
			argsToParse:    []string{"default", "default"},
			hasError:       false,
			checkFunc: func(app *App, exitCalled bool, exitRet int) error {
				expected := []string{"default"}
				if !reflect.DeepEqual(expected, app.ExtraArguments) {
					return fmt.Errorf("expected %v, got: %v", expected, app.ExtraArguments)
				}
				return nil
			},
		},
		{
			testName: "non-default subcommand with argument and value",
			appSubCommands: []*SubCommand{
				{Name: "default", Default: true},
				{Name: "extra"},
			},
			appFlags:    []*Flag{{Name: "known", TakesValue: true, Value: &SimpleValue{}}},
			argsToParse: []string{"--known", "value", "extra"},
			checkFunc: func(app *App, exitCalled bool, exitRet int) error {
				if app.FoundSubCommand.Name != "extra" {
					return fmt.Errorf("expected extra, got: %v", app.FoundSubCommand.Name)
				}
				return nil
			},
		},
		{
			testName: "preparse func fail",
			appFlags: []*Flag{{Name: "hello", PreParse: func(f *Flag, app *App) error {
				return errors.New("hello world")
			}}},
			hasError: true,
			errorStr: "hello world",
		},
		{
			testName: "preparse func succeed",
			appFlags: []*Flag{{Name: "hello", PreParse: func(f *Flag, app *App) error {
				return nil
			}}},
		},
		{
			testName: "onpresent func fail",
			appFlags: []*Flag{
				{
					Name: "hello",
					OnPresent: func(f *Flag, argString string, repeatCount int, app *App) error {
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
			appFlags: []*Flag{{Name: "hello", PreParse: func(f *Flag, app *App) error {
				return nil
			}}},
			argsToParse: []string{"--hello"},
		},
		{
			testName: "onpresent func not present",
			appFlags: []*Flag{
				{Name: "notused"},
				{Name: "hello", OnPresent: func(f *Flag, argString string, repeatCount int, app *App) error {
					return errors.New("hello world")
				}}},
			argsToParse: []string{"--notused"},
		},
		{
			testName:    "args name repeated",
			appFlags:    []*Flag{{Name: "v"}},
			argsToParse: []string{"-vvvvv"},
			checkFunc: func(app *App, exitCalled bool, exitRet int) error {
				arg := app.flagsMap["v"]
				if arg.RepeatCount != 5 {
					return fmt.Errorf("arg count mismatched, expected 5, got: %d", arg.RepeatCount)
				}
				return nil
			},
		},
		{
			testName: "onset 1",
			appFlags: []*Flag{{Name: "testonset", OnSet: func(f *Flag, argString string, newValue interface{}, app *App) error {
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
			appFlags: []*Flag{{Name: "testonset", OnSet: func(f *Flag, argString string, newValue interface{}, app *App) error {
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
			appFlags: []*Flag{{Name: "hello", PostParse: func(f *Flag, sc *SubCommand, app *App) error {
				return errors.New("hello world")
			}}},
			hasError: true,
			errorStr: "hello world",
		},
		{
			testName: "postparse func succeed",
			appFlags: []*Flag{{Name: "hello", PostParse: func(f *Flag, sc *SubCommand, app *App) error {
				return nil
			}}},
		},
		{
			testName: "postparse func ErrExitSuccess",
			appFlags: []*Flag{{Name: "hello", PostParse: func(f *Flag, sc *SubCommand, app *App) error {
				return ErrExitSuccess
			}}},
			checkFunc: func(app *App, exitCalled bool, exitRet int) error {
				if !exitCalled {
					return errors.New("no exit called")
				}
				if exitRet != 0 {
					return fmt.Errorf("exit ret expected 0, got: %d", exitRet)
				}
				return nil
			},
		},
	}

	for _, currentTest := range testList {
		t.Run(currentTest.testName, func(t *testing.T) {
			exitCalled := false
			exitRet := 0
			app := &App{
				exitFunc: func(i int) {
					exitCalled = true
					exitRet = i
				},
			}
			app.AddSubCommands(currentTest.appSubCommands...)
			app.AddFlags(currentTest.appFlags...)
			err := app.Parse(currentTest.argsToParse)
			if currentTest.hasError == (err == nil) {
				t.Fatalf("%q: expected error %v, got: %v", currentTest.testName, currentTest.hasError, err)
			}
			if err != nil && !strings.Contains(err.Error(), currentTest.errorStr) {
				t.Fatalf("test %q: expected %q in error: %v", currentTest.testName, currentTest.errorStr, err)
			}
			if currentTest.checkFunc != nil {
				if err := currentTest.checkFunc(app, exitCalled, exitRet); err != nil {
					t.Fatalf("test %q: check: %v", currentTest.testName, err)
				}
			}
		})
	}
}
