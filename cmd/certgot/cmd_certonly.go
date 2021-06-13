package main

import (
	"fmt"

	"github.com/eggsampler/certgot/cli"
)

var (
	cmdCertOnly = &cli.SubCommand{
		Name:       "certonly",
		Run:        commandCertOnly,
		HelpTopics: []string{TOPIC_COMMON},
	}
)

func commandCertOnly(app *cli.App) error {
	val := argAuthenticator.StringOrDefault()

	fmt.Println("authenticator: ", val)

	return nil
}
