package main

import (
	"github.com/eggsampler/certgot/cli"
)

const (
	CMD_CERTONLY = "certonly"
)

var (
	cmdCertOnly = &cli.Command{
		Name:           CMD_CERTONLY,
		RunFunc:        commandCertOnly,
		HelpCategories: []string{TOPIC_COMMON},
		HelpFlags:      []string{FLAG_NONINTERACTIVE},
	}
)

func commandCertOnly(ctx *cli.Context) error {
	/*
		val, err := flagAuthenticator.String(getFlagValues(true))
		if err != nil {
			return err
		}

		fmt.Println("authenticator: ", val)

	*/

	return nil
}
