package main

import (
	"github.com/eggsampler/certgot/cli"
)

const (
	CMD_HELP = "help"
)

var (
	cmdHelp = &cli.Command{
		Name:    CMD_HELP,
		RunFunc: commandHelp,
	}
)

func commandHelp(ctx *cli.Context) error {
	if len(ctx.ExtraArguments) > 0 {
		return ctx.App.PrintHelp(ctx, ctx.ExtraArguments[0])
	} else {
		return ctx.App.PrintHelp(ctx, "")
	}
}
