package main

import (
	"github.com/eggsampler/certgot/cli"
)

var (
	cmdRun = &cli.SubCommand{
		Name:       "run",
		Default:    true,
		HelpTopics: []string{TOPIC_COMMON},
		Usage: cli.SubCommandUsage{
			UsageDescription: "Obtain & install a certificate in your current webserver",
		},
		// TODO
	}
)
