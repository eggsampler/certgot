package main

import (
	"github.com/eggsampler/certgot/cli"
)

var (
	cmdRun = &cli.Command{
		Name:             "run",
		Default:          true,
		HelpCategories:   []string{TOPIC_COMMON},
		UsageDescription: "Obtain & install a certificate in your current webserver",
		// TODO
	}
)
