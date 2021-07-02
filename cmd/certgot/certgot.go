package main

import (
	"fmt"
	"os"

	"github.com/eggsampler/certgot/cli"
	"github.com/eggsampler/certgot/isdelve"
	"github.com/eggsampler/certgot/log"
)

func main() {
	app := &cli.App{
		Name: "certgot",

		Flags: cli.FlagList{
			flagHelp,
			flagConfig,
			flagWorkDir,
			flagLogsDir,
			flagConfigDir,
			flagDomains,
			flagCertName,
			flagNonInteractive,
			flagForceInteractive,
		},

		Commands: cli.CommandList{
			cmdRun,
			cmdCertOnly,
			cmdCertificates,
			cmdHelp,
		},

		Help: cli.HelpCategories{
			categoryUsage,
			categoryCommon,
			categoryManageCerts,
			categoryOptional,
		},
	}

	// pass through panics if running debug
	if !isdelve.Enabled {
		app.RecoverFunc = doRecover
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func doRecover(_ *cli.App, r interface{}) error {
	log.WithField("recover", r).Error("recovered")
	return nil
}
