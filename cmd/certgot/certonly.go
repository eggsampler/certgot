package main

import (
	"fmt"

	"github.com/eggsampler/certgot/cli"
)

func commandCertOnly(app *cli.App) error {
	val := argAuthenticator.StringOrDefault()

	fmt.Println("authenticator: ", val)

	return nil
}
