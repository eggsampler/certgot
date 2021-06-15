package main

import (
	"github.com/eggsampler/certgot/cli"
)

const (
	CMD_HELP = "help"
)

var (
	cmdHelp = &cli.SubCommand{
		Name: CMD_HELP,
		Run:  commandHelp,
	}
)

func commandHelp(app *cli.App) error {
	s := ""
	if len(app.ExtraArguments) > 0 {
		s = app.ExtraArguments[0]
	}
	app.PrintHelp(s)
	return nil
}
