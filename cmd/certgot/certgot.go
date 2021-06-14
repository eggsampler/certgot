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
	)

	app.AddArguments(
		argHelp,
		argConfig,
		argWorkDir,
		argLogsDir,
		argConfigDir,
		argEmail,
		argRegisterUnsafely,
		argStandAlone,
		argWebRoot,
		argAuthenticator,
		argDomains,
		argCertName)

	errExit := func(err error) {
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	}

	args := os.Args[1:]
	log.WithField("args", args).Debug("parsing arguments")
	errExit(app.Parse(args))
	log.WithField("config", argConfig.StringSliceOrDefault()).Debug("loading config")
	errExit(app.LoadConfig(argConfig))
	log.Debug("running")
	errExit(app.Run())
}

func doRecover(_ *cli.App, r interface{}) {
	log.WithField("recover", r).Error("recovered")
}
