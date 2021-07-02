package cli2

import "strings"

type CommandList []*Command

func (fl CommandList) Get(name string) *Command {
	for _, f := range fl {
		if strings.EqualFold(name, f.Name) {
			return f
		}
	}

	return nil
}

type Command struct {
	// Name is the name of the command and what is provided in the argument list to run this command
	Name string

	// Default is whether this command is the default command
	Default bool

	// HelpCategories is a list of names that this command should be printed in when printing help for a category
	HelpCategories []string

	// HelpFlags shows the list of flag arguments for the specific subcommand
	// Purely for help purposes
	HelpFlags []string

	// Usage is used to show custom usage args/params after `certgot subcommand XXXX`, where XXXX is Usage
	// If not supplied, by default shows `[options] ...`
	Usage string

	// UsageDescription is printed after the usage in help
	UsageDescription string

	// ArgumentDescription is description text shown before the Flags list
	ArgumentDescription string
}
