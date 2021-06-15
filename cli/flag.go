package cli

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

// RequireValueIfSet can be used in Flag.PostParse to ensure an argument, if present, has a value
// TODO: determine if this needs to be another property on Flag so that Flag.PostParse can still be used otherwise
func RequireValueIfSet() func(f *Flag, sc *SubCommand, app *App) error {
	return func(f *Flag, sc *SubCommand, app *App) error {
		if f.TakesValue && f.isPresent && !f.HasValue() {
			return errors.New("argument requires a value")
		}
		return nil
	}
}

type Flag struct {
	Name          string
	AltNames      []string
	DefaultValue  DefaultValue
	TakesValue    bool
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

	value interface{}
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

func (f Flag) HasValue() bool {
	return f.value != nil
}

func (f Flag) Value() interface{} {
	return f.value
}

func (f *Flag) Set(newValue interface{}) error {
	if !f.TakesValue {
		return fmt.Errorf("argument does not take a value")
	}
	if !f.TakesMultiple {
		f.value = newValue
		return nil
	}
	if f.value == nil {
		f.value = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(newValue)), 0, 10).Interface()
	}
	var valsToSet []reflect.Value
	if strVal, ok := newValue.(string); ok && strings.Contains(strVal, ",") {
		vals := strings.Split(strVal, ",")
		for _, sv := range vals {
			valsToSet = append(valsToSet, reflect.ValueOf(strings.TrimSpace(sv)))
		}
	} else {
		valsToSet = []reflect.Value{reflect.ValueOf(newValue)}
	}
	reflect.ValueOf(&f.value).Elem().Set(reflect.Append(reflect.ValueOf(f.value), valsToSet...))
	return nil
}

func (f Flag) String() string {
	s, _ := f.value.(string)
	return s
}

func (f Flag) StringOrDefault() string {
	if f.HasValue() {
		s, ok := f.value.(string)
		if ok {
			return s
		}
	}
	if f.DefaultValue == nil {
		return ""
	}
	s, _ := f.DefaultValue.Get().(string)
	return s
}

func (f Flag) StringSlice() []string {
	s, _ := f.value.([]string)
	return s
}

func (f Flag) StringSliceOrDefault() []string {
	if f.HasValue() {
		s, ok := f.value.([]string)
		if ok {
			return s
		}
	}
	if f.DefaultValue == nil {
		return nil
	}
	s, _ := f.DefaultValue.Get().([]string)
	return s
}

func (f Flag) Bool() bool {
	return f.value.(bool)
}

func (f Flag) BoolOrDefault() bool {
	if f.isPresent {
		b, _ := f.value.(bool)
		return b
	}
	if f.DefaultValue == nil {
		return false
	}
	defaultVal := f.DefaultValue.Get()
	if defaultVal == nil {
		return false
	}
	b, _ := defaultVal.(bool)
	return b
}
