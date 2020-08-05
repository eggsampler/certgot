package main

import (
	"github.com/eggsampler/certgot/constants"

	"github.com/eggsampler/certgot/cli"
)

var (
	argConfig = cli.Argument{
		Name:         constants.ARG_CONFIG,
		AltNames:     []string{"c"},
		DefaultValue: &cli.SimpleValue{Value: defaultConfigFiles},
		TakesValue:   true,
	}
	argWorkDir = cli.Argument{
		Name:         constants.ARG_WORKDIR,
		DefaultValue: &cli.SimpleValue{Value: defaultWorkDir},
		TakesValue:   true,
	}
	argLogsDir = cli.Argument{
		Name:         constants.ARG_LOGSDIR,
		DefaultValue: &cli.SimpleValue{Value: defaultLogsDir},
		TakesValue:   true,
	}
	argConfigDir = cli.Argument{
		Name:         constants.ARG_CONFIGDIR,
		DefaultValue: &cli.SimpleValue{Value: defaultConfigDir},
		TakesValue:   true,
	}

	argEmail = cli.Argument{
		Name: constants.ARG_EMAIL,
		DefaultValue: &cli.AskValue{
			Query:  "Enter email address (used for urgent renewal and security notices)",
			Cancel: "An e-mail address or --register-unsafely-without-email must be provided.",
		},
		TakesValue: true,
	}
	argRegisterUnsafely = cli.Argument{
		Name:         constants.ARG_REGISTER_UNSAFELY_WITHOUT_EMAIL,
		DefaultValue: &cli.SimpleValue{Value: false},
		TakesValue:   true,
	}

	argStandAlone = cli.Argument{
		Name:         constants.ARG_STANDALONE,
		DefaultValue: &cli.SimpleValue{Value: false},
	}
	argWebRoot = cli.Argument{
		Name:         constants.ARG_WEBROOT,
		DefaultValue: &cli.SimpleValue{Value: false},
	}
	argAuthenticator = cli.Argument{
		Name:     constants.ARG_AUTHENTICATOR,
		AltNames: []string{constants.ARG_AUTHENTICATOR_SHORT},
		DefaultValue: &cli.ListValue{
			Query:  "How would you like to authenticate with the ACME CA?",
			Cancel: "", // TODO: something here
			Options: []cli.ListValueOption{
				{
					Option: "Spin up a temporary webserver",
					Value:  "standalone",
				},
				{
					Option: "Place files in webroot directory",
					Value:  "webroot",
				},
			},
		},
	}

	helpTopics = []cli.HelpTopic{
		{
			Topics:          []string{"usage"},
			Name:            "usage",
			Description:     "certgot [SUBCOMMAND] [options] [-d DOMAIN] [-d DOMAIN] ...",
			LongDescription: "Certgot can obtain and install HTTPS/TLS/SSL certificates. By default, it will attempt to use a webserver both for obtaining and installing the certificate. The most common SUBCOMMANDS and flags are:",
			ShowFunc:        cli.ShowNotSubcommand,
		},
		{
			Topics:   []string{"common"},
			Name:     "obtain, install, and renew certificates",
			ShowFunc: cli.ShowNoTopic,
		},
	}
)

func setupSubCommands(app *cli.App) {
	app.AddSubCommands([]*cli.SubCommand{
		{
			Name: "certificates",
			Run:  commandCertificates,
		},
		{
			Name: "certonly",
			Run:  commandCertOnly,
		},
		{
			Name:    "run",
			Default: true,
			// TODO
		},
	}...)
}

func setupArguments(app *cli.App) {
	app.AddArguments(
		&argConfig,
		&argWorkDir,
		&argLogsDir,
		&argConfigDir,
		&argEmail,
		&argRegisterUnsafely,
		&argStandAlone,
		&argWebRoot,
		&argAuthenticator)
}

func setupHelp(app *cli.App) {
	app.AddHelpTopics(helpTopics...)
}
