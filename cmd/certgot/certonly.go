package main

import (
	"fmt"

	"github.com/eggsampler/certgot/cli"
)

func commandCertOnly(c *cli.Context) error {

	val := c.App.GetArgument("authenticator").StringOrDefault()

	fmt.Println("authenticator: ", val)

	return nil
}
