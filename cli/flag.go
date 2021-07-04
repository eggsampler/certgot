package cli

import (
	"fmt"
	"strconv"
	"strings"
)

type FlagList []*Flag

func (fl FlagList) Get(name string) *Flag {
	for _, f := range fl {
		if strings.EqualFold(name, f.Name) {
			return f
		}
		for _, an := range f.AltNames {
			if strings.EqualFold(an, f.Name) {
				return f
			}
		}
	}

	return nil
}

// Flag represents an argument that is prefixed by a single dash, or two dashes
type Flag struct {
	// Name is the name of the flag, ie the string after the dash(es)
	Name string

	// AltNames is other names the flag can go by
	// eg, can be used to add pluralised names, or short single character names
	AltNames []string

	// AllowShortRepeat determines whether short flags (single letter, single dash) are allowed to repeat
	// eg, -vvvvvvvv
	AllowShortRepeat bool

	// AllowMultiple determines whether the flag is allow to be present multiple times in the arguments
	// eg, --hello 1 --hello 2
	AllowMultiple bool

	// TakesValue determines whether the flag takes a value after the flag, or inline with the flag
	// eg, after the flag: `--hello world`, or inline: `--hello=world`
	TakesValue bool

	// OnSetFunc is a function that can be run after all the arguments have been parsed, but before a Command is run
	// eg, can be used to check some conditions based on the flag value
	OnSetFunc func(f *Flag, ctx *Context) error

	HelpDefault     func(ctx *Context) string
	HelpValueName   string
	HelpDescription string
	HelpCategories  []string

	flags  []string
	values []string
}

func (f Flag) Int() int {
	if len(f.values) == 0 {
		return 0
	}
	i, _ := strconv.Atoi(f.values[0])
	return i
}

func (f Flag) String() string {
	if len(f.values) == 0 {
		return ""
	}
	return f.values[0]
}

func (f Flag) StringSlice() []string {
	return f.values
}

func SetConfigValue(name string) func(f *Flag, ctx *Context) error {
	return func(f *Flag, ctx *Context) error {
		cfg := ctx.App.Configs.Get(name)
		if cfg == nil {
			return fmt.Errorf("no config %q for flag %q", name, f.Name)
		}
		cfg.set(f.values)
		return nil
	}
}
