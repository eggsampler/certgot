package main

import (
	"github.com/eggsampler/certgot/cli"
)

var (
	argHelp = &cli.Argument{
		Name:       ARG_HELP,
		AltNames:   []string{ARG_HELP_SHORT},
		TakesValue: true,
		HelpTopics: []string{TOPIC_OPTIONAL},
		Usage: cli.ArgumentUsage{
			Description: "show this help message and exit",
		},
		PostParse: func(arg *cli.Argument, sc *cli.SubCommand, app *cli.App) error {
			if !arg.IsPresent {
				return nil
			}
			app.PrintHelp()
			return cli.ErrExitSuccess
		},
	}
	argConfig = &cli.Argument{
		Name:         ARG_CONFIG,
		AltNames:     []string{ARG_CONFIG_SHORT},
		DefaultValue: &cli.SimpleValue{Value: defaultConfigFiles},
		TakesValue:   true,
		HelpTopics:   []string{TOPIC_OPTIONAL},
		Usage: cli.ArgumentUsage{
			ArgName:     "CONFIG_FILE",
			Description: "path to config file",
		},
	}
	argWorkDir = &cli.Argument{
		Name:         ARG_WORKDIR,
		DefaultValue: &cli.SimpleValue{Value: defaultWorkDir},
		TakesValue:   true,
		HelpTopics:   []string{TOPIC_OPTIONAL},
	}
	argLogsDir = &cli.Argument{
		Name:         ARG_LOGSDIR,
		DefaultValue: &cli.SimpleValue{Value: defaultLogsDir},
		TakesValue:   true,
		HelpTopics:   []string{TOPIC_OPTIONAL},
	}
	argConfigDir = &cli.Argument{
		Name:         ARG_CONFIGDIR,
		DefaultValue: &cli.SimpleValue{Value: defaultConfigDir},
		TakesValue:   true,
		HelpTopics:   []string{TOPIC_OPTIONAL},
	}

	argEmail = &cli.Argument{
		Name: ARG_EMAIL,
		DefaultValue: &cli.AskValue{
			Query:  "Enter email address (used for urgent renewal and security notices)",
			Cancel: "An e-mail address or --register-unsafely-without-email must be provided.",
		},
		TakesValue: true,
	}
	argRegisterUnsafely = &cli.Argument{
		Name:         ARG_REGISTER_UNSAFELY_WITHOUT_EMAIL,
		DefaultValue: &cli.SimpleValue{Value: false},
		TakesValue:   true,
	}

	argStandAlone = &cli.Argument{
		Name:         ARG_STANDALONE,
		DefaultValue: &cli.SimpleValue{Value: false},
	}
	argWebRoot = &cli.Argument{
		Name:         ARG_WEBROOT,
		DefaultValue: &cli.SimpleValue{Value: false},
	}
	argAuthenticator = &cli.Argument{
		Name:     ARG_AUTHENTICATOR,
		AltNames: []string{ARG_AUTHENTICATOR_SHORT},
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
			Topic:           TOPIC_USAGE,
			Name:            "usage",
			Description:     "\n  certgot [SUBCOMMAND] [options] [-d DOMAIN] [-d DOMAIN] ...",
			LongDescription: "Certgot can obtain and install HTTPS/TLS/SSL certificates. By default, it will attempt to use a webserver both for obtaining and installing the certificate. The most common SUBCOMMANDS and flags are:",
			ShowFunc:        cli.ShowNotSubcommand,
		},
		{
			Topic:    TOPIC_COMMON,
			Name:     "obtain, install, and renew certificates",
			ShowFunc: cli.ShowNoTopic,
		},
		{
			Topic:    TOPIC_MANAGE,
			Name:     "manage certificates",
			ShowFunc: cli.ShowNoTopic,
		},
		{
			Topic:    TOPIC_OPTIONAL,
			Name:     "optional arguments",
			ShowFunc: cli.ShowAlways,
		},
	}
)

func setupSubCommands(app *cli.App) {
	app.AddSubCommands(
		cmdRun,
		cmdCertOnly,
		cmdCertificates,
	)
}

func setupArguments(app *cli.App) {
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
		argAuthenticator)
}

func setupHelp(app *cli.App) {
	app.AddHelpTopics(helpTopics...)
}
