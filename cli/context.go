package cli

// Context represents a running context for a cli application, and is passed around to be used
type Context struct {

	// App is just a reference to the defined app, so you don't need to store a reference
	App *App

	// RawArguments is what is provided to App.Run
	RawArguments []string

	// Flags is a list of flags that is present after parsing
	// NB: These are only flags present in the arguments, all flags are still available in App.Flags
	Flags FlagList

	// Command is the specified (or default command if no command was specified)
	Command *Command

	// Just holds whether the default command was used
	IsDefaultCommand bool

	// Any extra arguments after a command is found and not belonging to any flags
	ExtraArguments []string
}
