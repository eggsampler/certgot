package main

import (
	"fmt"
	"os"

	"github.com/eggsampler/certgot/cli"
	"github.com/eggsampler/certgot/isdelve"
	"github.com/eggsampler/certgot/log"
)

func init() {
	log.SetLevel(log.MinLevel)
}

func main() {
	app := &cli.App{
		Name: "certgot",

		Flags: cli.FlagList{
			flagHelp,
			flagConfigFile,
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

		Configs: cli.ConfigList{
			cfgConfigFile,
		},

		Help: cli.HelpCategories{
			categoryUsage,
			categoryCommon,
			categoryManageCerts,
			categoryOptional,
		},

		PreRunFunc: func(ctx *cli.Context) error {
			cfg := ctx.App.Configs.Get(CONFIG_FILE)
			return ctx.App.LoadConfig(cfg.StringSlice(), !cfg.IsSet())
		},
	}

	// pass through panics if running debug
	if !isdelve.Enabled {
		app.RecoverFunc = doRecover
	}

	log.WithField("args", os.Args).Debug("running")

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
