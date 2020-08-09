package main

import (
	"fmt"

	"github.com/eggsampler/certgot/constants"

	"github.com/eggsampler/certgot/cli"
)

var (
	cmdCertonly = cli.SubCommand{
		Name:       "certonly",
		Run:        commandCertOnly,
		HelpTopics: []string{constants.TOPIC_COMMON},
	}
)

func commandCertOnly(app *cli.App) error {
	val := argAuthenticator.StringOrDefault()

	fmt.Println("authenticator: ", val)

	return nil
}
