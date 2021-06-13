package cli

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ErrExitSuccess represents an error that can be returned from an argument Argument.PostParse func
// if returned, program exits normally, error return 0
// TODO: use this for other functions on Argument ?
var ErrExitSuccess = errors.New("done")

type Argument struct {
	Name          string
	AltNames      []string
	DefaultValue  Value
	TakesValue    bool
	RepeatCount   int
	TakesMultiple bool

	HelpTopics []string
	Usage      ArgumentUsage

	PreParse  func(arg *Argument, app *App) error
	OnPresent func(arg *Argument, argString string, repeatCount int, app *App) error
	OnSet     func(arg *Argument, argString string, newValue interface{}, app *App) error
	PostParse func(arg *Argument, sc *SubCommand, app *App) error

	IsPresent bool

	value interface{}
}

type ArgumentUsage struct {
	ArgName     string
	Description string
}

func (arg Argument) HasValue() bool {
	return arg.value != nil
}

func (arg Argument) Value() interface{} {
	return arg.value
}

func (arg *Argument) Set(newValue interface{}) error {
	if !arg.TakesValue {
		return fmt.Errorf("argument does not take a value")
	}
	if !arg.TakesMultiple {
		arg.value = newValue
		return nil
	}
	if arg.value == nil {
		arg.value = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(newValue)), 0, 10).Interface()
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
	reflect.ValueOf(&arg.value).Elem().Set(reflect.Append(reflect.ValueOf(arg.value), valsToSet...))
	return nil
}

func (arg Argument) String() string {
	s, _ := arg.value.(string)
	return s
}

func (arg Argument) StringOrDefault() string {
	if arg.HasValue() {
		s, ok := arg.value.(string)
		if ok {
			return s
		}
	}
	if arg.DefaultValue == nil {
		return ""
	}
	s, _ := arg.DefaultValue.Get().(string)
	return s
}

func (arg Argument) StringSlice() []string {
	s, _ := arg.value.([]string)
	return s
}

func (arg Argument) Bool() bool {
	return arg.value.(bool)
}

func (arg Argument) BoolOrDefault() bool {
	if arg.IsPresent {
		b, _ := arg.value.(bool)
		return b
	}
	if arg.DefaultValue == nil {
		return false
	}
	defaultVal := arg.DefaultValue.Get()
	if defaultVal == nil {
		return false
	}
	b, _ := defaultVal.(bool)
	return b
}

func (arg Argument) StringSliceOrDefault() []string {
	if arg.HasValue() {
		s, ok := arg.value.([]string)
		if ok {
			return s
		}
	}
	if arg.DefaultValue == nil {
		return nil
	}
	s, _ := arg.DefaultValue.Get().([]string)
	return s
}
