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
		Name:        "certgot",
		FuncPreRun:  doPreRun,
		FuncPostRun: doPostRun,
	}

	// pass through panics if running debug
	if !isdelve.Enabled {
		app.FuncRecover = doRecover
	}

	app.AddHelpTopics(helpTopics...)

	app.AddSubCommands(
		cmdRun,
		cmdCertOnly,
		cmdCertificates,
		cmdHelp,
	)

	app.AddFlags(
		flagHelp,
		flagConfig,
		flagWorkDir,
		flagLogsDir,
		flagConfigDir,
		flagEmail,
		flagRegisterUnsafely,
		flagStandAlone,
		flagWebRoot,
		flagAuthenticator,
		flagDomains,
		flagCertName)

	errExit := func(err error) {
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	}

	args := os.Args[1:]
	log.WithField("flags", args).Debug("parsing arguments")
	errExit(app.Parse(args))
	log.WithField("config", flagConfig.StringSliceOrDefault()).Debug("loading config")
	errExit(app.LoadConfig(flagConfig))
	log.Debug("running")
	errExit(app.Run())
}

func doRecover(_ *cli.App, r interface{}) {
	log.WithField("recover", r).Error("recovered")
}
