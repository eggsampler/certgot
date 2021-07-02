package cli2

type Context struct {
	App            *App
	RawArguments   []string
	Flags          FlagList
	Command        *Command
	ExtraArguments []string
}
