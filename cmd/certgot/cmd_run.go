package main

import (
	"github.com/eggsampler/certgot/cli"
	"github.com/eggsampler/certgot/constants"
)

var (
	cmdRun = cli.SubCommand{
		Name:       "run",
		Default:    true,
		HelpTopics: []string{constants.TOPIC_COMMON},
		Usage: cli.Usage{
			Description: "Obtain & install a certificate in your current webserver",
		},
		// TODO
	}
)
