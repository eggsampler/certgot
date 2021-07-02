package cli_old

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var ErrNoValueSet = errors.New("no value set")

// Value represents a type used for a Flag.DefaultValue
// These can be used to represent simple types like a string, or boolean
// But also more complex types, for example ones that can prompt user for input
type Value interface {
	// Get returns the value set
	Get(nonInteractive, forceInteractive, isTerminal, includeDefault bool) (interface{}, error)

	// Set sets the value
	Set(interface{}) Value

	// IsSet returns if the value has been Set
	IsSet() bool

	// UsageDefault shows a default value when printing help
	// Used to show a long string describing the default behaviour
	UsageDefault() string

	// HelpDefault shows a default value when printing help
	// Used to show the default value, if set
	HelpDefault() string
}

func NewSimpleValueDefault(value interface{}) Value {
	return &SimpleValue{
		value: value,
		set:   true,
	}
}

func NewSimpleValueNotSet(usageDefault string) Value {
	return &SimpleValue{
		usageDefault: usageDefault,
	}
}

type SimpleValue struct {
	value        interface{}
	def          interface{}
	set          bool
	usageDefault string
}

func (sv SimpleValue) Get(nonInteractive, forceInteractive, isTerminal, includeDefault bool) (interface{}, error) {
	if !sv.set {
		if includeDefault && sv.def != nil {
			return sv.def, nil
		}
		return nil, ErrNoValueSet
	}
	return sv.value, nil
}

func (sv *SimpleValue) Set(v interface{}) Value {
	sv.set = true
	sv.value = v
	return sv
}

func (sv SimpleValue) IsSet() bool {
	return sv.set
}

func (sv SimpleValue) HelpDefault() string {
	return fmt.Sprintf("%+v", sv.def)
}

func (sv SimpleValue) UsageDefault() string {
	if sv.usageDefault != "" {
		return sv.usageDefault
	}
	if sv.value != nil {
		// TODO: should this be some variant of %v ?
		return fmt.Sprintf("%s", sv.value)
	}
	return ""
}

type AskValue struct {
	SimpleValue
	Query  string
	Cancel string
}

func (av *AskValue) Get(nonInteractive, forceInteractive, isTerminal, includeDefault bool) (interface{}, error) {
	// if running in non-interactive mode, or there is no terminal output
	if nonInteractive || !isTerminal {
		return av.Get(nonInteractive, forceInteractive, isTerminal, includeDefault)
	}
	// if forcing interactive mode, or running in a terminal, and no value is set
	if (forceInteractive || isTerminal) && !av.IsSet() {
		// query the user for the value
		fmt.Printf("%s (Enter 'c' to cancel): ", av.Query)
		av.Set(readLine(av.Cancel))
	}

	return av.Get(nonInteractive, forceInteractive, isTerminal, includeDefault)
}

func (av AskValue) UsageDefault() string {
	return "Ask"
}

type ListValueOption struct {
	Option string
	Value  interface{}
}

type ListValue struct {
	SimpleValue
	Query   string
	Cancel  string
	Options []ListValueOption
}

func (lv *ListValue) Get(nonInteractive, forceInteractive, isTerminal, includeDefault bool) (interface{}, error) {
	// if running in non-interactive mode, or there is no terminal output
	if nonInteractive || !isTerminal {
		return lv.Get(nonInteractive, forceInteractive, isTerminal, includeDefault)
	}
	// if forcing interactive mode, or running in a terminal, and no value is set
	if (forceInteractive || isTerminal) && !lv.IsSet() {
		fmt.Println(lv.Query)
		fmt.Println(strings.Repeat("- ", 40))
		for k, v := range lv.Options {
			fmt.Printf("%d: %s (%v)\n", k+1, v.Option, v.Value)
		}
		fmt.Println(strings.Repeat("- ", 40))
		fmt.Printf("Select the appropriate number [1-%d] then [enter] (press 'c' to cancel): ", len(lv.Options))
		num, _ := strconv.Atoi(readLine(lv.Cancel))
		if num <= 0 || num > len(lv.Options) {
			fmt.Println("Invalid value")
			os.Exit(1)
		}
		lv.Set(lv.Options[num-1].Value)
	}

	return lv.Get(nonInteractive, forceInteractive, isTerminal, includeDefault)
}

func (lv ListValue) UsageDefault() string {
	return "Ask"
}

func readLine(cancel string) string {
	scanner := bufio.NewScanner(os.Stdin)
	success := scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
		return ""
	}
	input := scanner.Text()
	if !success || strings.ToLower(input) == "c" {
		if cancel != "" {
			fmt.Println(cancel)
		}
		os.Exit(1)
		return ""
	}
	return input
}
