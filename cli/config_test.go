package cli

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_parseConfig(t *testing.T) {
	emptyCfg := func(cfg []configEntry) error {
		if len(cfg) != 0 {
			return fmt.Errorf("unexpected config: %v", cfg)
		}
		return nil
	}

	testList := []struct {
		testName     string
		configString string
		fileName     string
		hasError     bool
		errorStr     string
		checkFunc    func(cfg []configEntry) error
	}{
		{
			testName:  "empty config",
			checkFunc: emptyCfg,
		},
		{
			testName:     "invalid",
			configString: "1234",
			hasError:     true,
			errorStr:     "invalid",
		},
		{
			testName:     "no config",
			configString: "# hello world",
			checkFunc:    emptyCfg,
		},
		{
			testName:     "config w/ value",
			configString: "hello=world",
			fileName:     "hi2u",
			checkFunc: func(cfg []configEntry) error {
				if len(cfg) != 1 {
					return errors.New("not 1 config")
				}
				entry := cfg[0]
				otherEntry := configEntry{
					fileName: "hi2u",
					line:     1,
					key:      "hello",
					hasValue: true,
					value:    "world",
				}
				if !reflect.DeepEqual(entry, otherEntry) {
					return fmt.Errorf("entry mismatch %+v != %+v", entry, otherEntry)
				}
				return nil
			},
		},
		{
			testName:     "config w/o value",
			configString: "hello",
			fileName:     "hi2u2",
			checkFunc: func(cfg []configEntry) error {
				if len(cfg) != 1 {
					return errors.New("not 1 config")
				}
				entry := cfg[0]
				otherEntry := configEntry{
					fileName: "hi2u2",
					line:     1,
					key:      "hello",
				}
				if !reflect.DeepEqual(entry, otherEntry) {
					return fmt.Errorf("entry mismatch %+v != %+v", entry, otherEntry)
				}
				return nil
			},
		},
	}

	for _, currentTest := range testList {
		cfg, err := parseConfig(strings.NewReader(currentTest.configString), currentTest.fileName)
		if currentTest.hasError == (err == nil) {
			t.Fatalf("%q: expected error %v, got: %v", currentTest.testName, currentTest.hasError, err)
		}
		if err != nil && !strings.Contains(err.Error(), currentTest.errorStr) {
			t.Fatalf("test %q: expected %q in error: %v", currentTest.testName, currentTest.errorStr, err)
		}
		if currentTest.checkFunc != nil {
			if err := currentTest.checkFunc(cfg); err != nil {
				t.Fatalf("%q: unexpected error %v", currentTest.testName, err)
			}
		}
	}
}

func Test_setConfig(t *testing.T) {
	testList := []struct {
		testName  string
		config    []configEntry
		args      map[string]*Argument
		hasError  bool
		errorStr  string
		checkFunc func(config map[string]configEntry, args map[string]*Argument) error
	}{
		{testName: "empty"},
		{
			testName: "no arg for cfg",
			config:   []configEntry{{}},
			hasError: true,
			errorStr: "unknown argument",
		},
		{
			testName: "set arg",
			config:   []configEntry{{key: "hello", hasValue: true, value: "world"}},
			args: map[string]*Argument{"hello": {
				TakesValue: true,
			}},
			checkFunc: func(config map[string]configEntry, args map[string]*Argument) error {
				arg := args["hello"]
				val := arg.String()
				if val != "world" {
					return fmt.Errorf("world != %s", val)
				}
				return nil
			},
		},
		{
			testName: "set arg fail",
			config:   []configEntry{{key: "hello", hasValue: true, value: "world"}},
			args:     map[string]*Argument{"hello": {}},
			hasError: true,
			errorStr: "error setting arg",
		},
	}

	for _, currentTest := range testList {
		err := setConfig(currentTest.config, currentTest.args)
		if currentTest.hasError == (err == nil) {
			t.Fatalf("%q: expected error %v, got: %v", currentTest.testName, currentTest.hasError, err)
		}
		if err != nil && !strings.Contains(err.Error(), currentTest.errorStr) {
			t.Fatalf("test %q: expected %q in error: %v", currentTest.testName, currentTest.errorStr, err)
		}
	}
}
