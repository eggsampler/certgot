package cli

type SubCommand struct {
	Name       string
	Default    bool
	Run        func(app *App) error
	HelpTopics []string
	Usage      SubCommandUsage
}

type SubCommandUsage struct {
	// LongUsage is used to show custom usage args/params after `certgot subcommand XXXX`, where XXXX is LongUsage
	// If not supplied, by default shows `[options] ...`
	LongUsage string

	// UsageDescription is printed after the usage in help
	UsageDescription string

	// ArgumentDescription is description text shown before the Flags list
	ArgumentDescription string

	// Flags shows the list of flag arguments for the specific subcommand
	// Purely for help purposes
	Flags []string
}
