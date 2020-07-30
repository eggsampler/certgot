package main

import (
	"path/filepath"

	"github.com/eggsampler/certgot/cli"
	"github.com/eggsampler/certgot/constants"
)

var (
	configFiles = []string{
		filepath.Join(string(filepath.Separator), "etc", "letsencrypt", "cli.ini"),
		filepath.Join("~", ".config", "letsencrypt", "cli.ini"),
	}
)

func setupApp(app *cli.App) {
	setupSubCommands(app)
	setupArguments(app)
}

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
	// TODO: os specific paths
	app.AddArguments([]*cli.Argument{
		{
			Name:         constants.ARG_CONFIG,
			AltNames:     []string{"c"},
			DefaultValue: &cli.SimpleValue{Value: configFiles},
			TakesValue:   true,
		},

		{
			Name:         constants.ARG_WORKDIR,
			DefaultValue: &cli.SimpleValue{Value: filepath.Join("var", "lib", "letsencrypt")},
			TakesValue:   true,
		},
		{
			Name:         constants.ARG_LOGSDIR,
			DefaultValue: &cli.SimpleValue{Value: filepath.Join("var", "logs", "letsencrypt")},
			TakesValue:   true,
		},
		{
			Name:         constants.ARG_CONFIGDIR,
			DefaultValue: &cli.SimpleValue{filepath.Join("etc", "letsencrypt")},
			TakesValue:   true,
		},

		{
			Name: "email",
			DefaultValue: &cli.AskValue{
				Query:  "Enter email address (used for urgent renewal and security notices)",
				Cancel: "An e-mail address or --register-unsafely-without-email must be provided.",
			},
			TakesValue: true,
		},
		{
			Name:         "register-unsafely-without-email",
			DefaultValue: &cli.SimpleValue{false},
			TakesValue:   true,
		},

		{
			Name:         "standalone",
			DefaultValue: &cli.SimpleValue{false},
		},
		{
			Name:         "webroot",
			DefaultValue: &cli.SimpleValue{false},
		},
		{
			Name:     "authenticator",
			AltNames: []string{"a"},
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
		},
	}...)
}
