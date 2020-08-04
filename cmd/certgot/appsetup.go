package main

import (
	"path/filepath"

	"github.com/eggsampler/certgot/constants"

	"github.com/eggsampler/certgot/cli"
)

// TODO: os specific paths

var (
	argConfig = cli.Argument{
		Name:     constants.ARG_CONFIG,
		AltNames: []string{"c"},
		DefaultValue: &cli.SimpleValue{Value: []string{
			filepath.Join(string(filepath.Separator), "etc", "letsencrypt", "cli.ini"),
			filepath.Join("~", ".config", "letsencrypt", "cli.ini"),
		}},
		TakesValue: true,
	}
	argWorkDir = cli.Argument{
		Name:         constants.ARG_WORKDIR,
		DefaultValue: &cli.SimpleValue{Value: filepath.Join("var", "lib", "letsencrypt")},
		TakesValue:   true,
	}
	argLogsDir = cli.Argument{
		Name:         constants.ARG_LOGSDIR,
		DefaultValue: &cli.SimpleValue{Value: filepath.Join("var", "logs", "letsencrypt")},
		TakesValue:   true,
	}
	argConfigDir = cli.Argument{
		Name:         constants.ARG_CONFIGDIR,
		DefaultValue: &cli.SimpleValue{Value: filepath.Join("etc", "letsencrypt")},
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
