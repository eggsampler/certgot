package main

import "github.com/eggsampler/certgot/cli"

const (
	TOPIC_USAGE               = "usage"
	TOPIC_COMMON              = "common"
	TOPIC_MANAGE_CERTIFICATES = "manage"
	TOPIC_OPTIONAL            = "optional"
)

var (
	categoryUsage = &cli.HelpCategory{
		Category:         TOPIC_USAGE,
		Name:             "usage",
		Usage:            "certgot [SUBCOMMAND] [options] [-d DOMAIN] [-d DOMAIN] ...",
		UsageDescription: "Certgot can obtain and install HTTPS/TLS/SSL certificates. By default, it will attempt to use a webserver both for obtaining and installing the certificate. The most common SUBCOMMANDS and flags are:",
		ShowFunc:         cli.ShowNoCommand,
	}
	categoryCommon = &cli.HelpCategory{
		Category: TOPIC_COMMON,
		Name:     "obtain, install, and renew certificates",
		ShowFunc: cli.ShowNoCategory,
	}
	categoryManageCerts = &cli.HelpCategory{
		Category: TOPIC_MANAGE_CERTIFICATES,
		Name:     "manage certificates",
		ShowFunc: cli.ShowNoCategory,
	}
	categoryOptional = &cli.HelpCategory{
		Category: TOPIC_OPTIONAL,
		Name:     "optional arguments",
		ShowFunc: cli.ShowAlways,
	}
)
