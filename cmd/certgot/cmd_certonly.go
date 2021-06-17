package main

import (
	"fmt"

	"github.com/eggsampler/certgot/cli"
)

const (
	CMD_CERTONLY = "certonly"
)

var (
	cmdCertOnly = &cli.SubCommand{
		Name:       CMD_CERTONLY,
		Run:        commandCertOnly,
		HelpTopics: []string{TOPIC_COMMON},
		Flags:      []string{FLAG_NONINTERACTIVE},
	}
)

func commandCertOnly(app *cli.App) error {
	val, err := flagAuthenticator.String(getFlagValues(true))
	if err != nil {
		return err
	}

	fmt.Println("authenticator: ", val)

	return nil
}
