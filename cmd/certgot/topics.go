package main

import "github.com/eggsampler/certgot/cli"

const (
	TOPIC_USAGE               = "usage"
	TOPIC_COMMON              = "common"
	TOPIC_MANAGE_CERTIFICATES = "manage"
	TOPIC_OPTIONAL            = "optional"
)

var (
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
			Topic:    TOPIC_MANAGE_CERTIFICATES,
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
