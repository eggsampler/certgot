package main

import (
	"fmt"
	"os"

	"github.com/eggsampler/certgot/cli"
	"github.com/eggsampler/certgot/log"
	"github.com/eggsampler/certgot/util"
)

func main() {
	app := &cli.App{
		FuncPreRun:  doPreRun,
		FuncPostRun: doPostRun,
		FuncRecover: doRecover,
	}

	setupApp(app)

	if err := app.LoadConfig(configFiles); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	err := app.Run(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func doPreRun(c *cli.Context) error {
	ll := log.Entry{}
	for _, v := range c.App.GetArguments() {
		if !v.HasValue() {
			continue
		}
		ll = ll.WithField(v.Name, v.Value())
	}
	ll.Info("cli arguments")
	log.WithField("subcommand", c.SubCommand.Name).Info("cli subcommand")

	dirs := c.GetDirectories()

	if err := util.SetupDirectories(dirs); err != nil {
		return fmt.Errorf("Error setting up directories: %v\n", err)
	}

	if err := util.SetupLocks(dirs); err != nil {
		return fmt.Errorf("Error setting up lock files: %v\n", err)
	}

	return nil
}

func doPostRun(c *cli.Context) error {
	util.CleanupLocks()

	return nil
}

func doRecover(c *cli.Context, r interface{}) error {
	log.WithField("recover", r).Error("recovered")

	return nil
}
