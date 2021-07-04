package cli

import (
	"strconv"
	"strings"
)

// This is not a map because a flag can be looked up by it's alternative names
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

func (fl *FlagList) Put(f *Flag) {
	existing := fl.Get(f.Name)
	if existing != nil {
		return
	}
	*fl = append(*fl, f)
}

type FlagValue struct {
	FlagName string
	RawFlag  string
	HasValue bool
	Value    string
}

// Flag represents an argument that is prefixed by a single dash, or two dashes
type Flag struct {
	// Name is the name of the flag, ie the string after the dash(es)
	Name string

	// AltNames is other names the flag can go by
	// eg, can be used to add pluralised names, or short single character names
	AltNames []string

	// AllowShortRepeat determines whether short flags (single dash, single letter) are allowed to repeat
	// eg, -vvvvvvvv
	AllowShortRepeat bool

	// AllowMultiple determines whether the flag is allow to be present multiple times in the arguments
	// eg, --hello 1 --hello 2
	AllowMultiple bool

	// TakesValue determines whether the flag takes a value after the flag, or inline with the flag
	// eg, after the flag: `--hello world`, or inline: `--hello=world`
	TakesValue    bool
	RequiresValue bool

	// PostParseFunc is a function that can be run after all the arguments have been parsed, but before a Command is run
	// eg, can be used to check some conditions based on the flag value
	PostParseFunc func(f *Flag, ctx *Context) error

	// HelpDefault returns a string to show in the help list in brackets for the default value
	HelpDefault func(ctx *Context) (string, error)

	// HelpValueName shows text after the flag in the help list (eg, `--blah WORLD` for HelpValueName = WORLD)
	HelpValueName string

	// HelpDescription is some descriptive help text for the flag
	HelpDescription string

	// HelpCategories identifies which categories to show the flag under
	HelpCategories []string

	// valuesInfo contains extra info about the flag, ie. the raw flag so short flags can be counted etc
	valuesInfo []FlagValue

	// valuesRaw is just a list of the values that is set for the flag
	valuesRaw []string
}

func (f Flag) ValueList() []FlagValue {
	return f.valuesInfo
}

func (f Flag) Int() int {
	if len(f.valuesRaw) == 0 {
		return 0
	}
	i, _ := strconv.Atoi(f.valuesRaw[0])
	return i
}

func (f Flag) String() string {
	if len(f.valuesRaw) == 0 {
		return ""
	}
	return f.valuesRaw[0]
}

func (f Flag) StringSlice() []string {
	return f.valuesRaw
}
