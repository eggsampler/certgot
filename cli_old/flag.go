package cli_old

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ErrExitSuccess represents an error that can be returned from an argument Flag.PostParse func
// if returned, program exits normally, error return 0
// TODO: use this for other functions on Flag ?
var ErrExitSuccess = errors.New("done")

var ErrRequiresValue = errors.New("argument requires a value")

// TODO: this error seems hacky, definitely rethink value
var ErrNoValueProvided = errors.New("no value provided")

// RequireValueIfSet can be used in Flag.PostParse to ensure an argument, if present, has a value
// TODO: determine if this needs to be another property on Flag so that Flag.PostParse can still be used otherwise
func RequireValueIfSet() func(f *Flag, sc *SubCommand, app *App) error {
	return func(f *Flag, sc *SubCommand, app *App) error {
		if f.TakesValue && f.isPresent && !f.Value.IsSet() {
			return ErrRequiresValue
		}
		return nil
	}
}

type Flag struct {
	Name          string
	AltNames      []string
	TakesValue    bool
	Value         Value // TODO: not expose this field somehow?
	RepeatCount   int
	TakesMultiple bool

	HelpTopics []string
	Usage      ArgumentUsage

	PreParse  func(f *Flag, app *App) error
	OnPresent func(f *Flag, argString string, repeatCount int, app *App) error
	OnSet     func(f *Flag, argString string, newValue interface{}, app *App) error
	PostParse func(f *Flag, sc *SubCommand, app *App) error

	isPresent           bool
	isPresentInArgument bool
}

type ArgumentUsage struct {
	ArgName     string
	Description string
}

func (f Flag) IsPresent() bool {
	return f.isPresent
}

func (f Flag) IsPresentInArgument() bool {
	return f.isPresentInArgument
}

func (f *Flag) Set(newValue interface{}) error {
	if f.Value == nil {
		return ErrNoValueProvided
	}
	// throw an error if this flag doesn't explicitly take a value
	if !f.TakesValue {
		return fmt.Errorf("argument does not take a value")
	}
	// set the value if the flag doesn't take multiples (multiples being: -d blah1 -d blah2, or -d blah1,blah2)
	if !f.TakesMultiple {
		f.Value.Set(newValue)
		return nil
	}
	// the variable to store the value in for multiple types
	var v interface{}
	if !f.Value.IsSet() {
		// first time value has been set (ie, for first instance of flag being parsed)
		v = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(newValue)), 0, 10).Interface()
	} else {
		// flag has already been set once, grab the underlying value to use to append more values to
		var err error
		v, err = f.Value.Get(true, false, false, false)
		if err != nil {
			return fmt.Errorf("error getting underlying value to append multiples: %w", err)
		}
	}
	// if the new value is a string, and contains multiple comma-separated elements, extract each one
	var valsToSet []reflect.Value
	if strVal, ok := newValue.(string); ok && strings.Contains(strVal, ",") {
		vals := strings.Split(strVal, ",")
		for _, sv := range vals {
			valsToSet = append(valsToSet, reflect.ValueOf(strings.TrimSpace(sv)))
		}
	} else {
		// otherwise, just use the single value
		valsToSet = []reflect.Value{reflect.ValueOf(newValue)}
	}
	// append parsed & extracted values to underlying value
	reflect.ValueOf(&v).Elem().Set(reflect.Append(reflect.ValueOf(v), valsToSet...))
	// and store to underlying value itself
	f.Value.Set(v)
	return nil
}

func (f Flag) String(nonInteractive, forceInteractive, isTerminal, includeDefault bool) (string, error) {
	if f.Value == nil {
		return "", ErrNoValueProvided
	}
	v, err := f.Value.Get(nonInteractive, forceInteractive, isTerminal, includeDefault)
	if err != nil {
		return "", err
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("value type is not a string (got: %T)", v)
	}
	return s, nil
}

func (f Flag) StringSlice(nonInteractive, forceInteractive, isTerminal, includeDefault bool) ([]string, error) {
	if f.Value == nil {
		return nil, ErrNoValueProvided
	}
	v, err := f.Value.Get(nonInteractive, forceInteractive, isTerminal, includeDefault)
	if err != nil {
		return nil, err
	}
	s, ok := v.([]string)
	if !ok {
		return nil, fmt.Errorf("value type is not a string slice (got: %T)", v)
	}
	return s, nil
}

func (f Flag) Bool(nonInteractive, forceInteractive, isTerminal, includeDefault bool) (bool, error) {
	if f.Value == nil {
		return false, ErrNoValueProvided
	}
	v, err := f.Value.Get(nonInteractive, forceInteractive, isTerminal, includeDefault)
	if err != nil {
		return false, err
	}
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("value type is not a bool (got: %T)", v)
	}
	return b, nil
}
