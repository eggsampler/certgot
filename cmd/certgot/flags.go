package main

import "github.com/eggsampler/certgot/cli"

// TODO: pick a better naming scheme to identify the constant names vs the variable flags
// to better match go naming https://golang.org/doc/effective_go#mixed-caps

const (
	FLAG_HELP                            = "help"
	FLAG_HELP_SHORT                      = "h"
	FLAG_CONFIG                          = "config"
	FLAG_CONFIG_SHORT                    = "c"
	FLAG_WORK_DIR                        = "work-dir"
	FLAG_LOGS_DIR                        = "logs-dir"
	FLAG_CONFIG_DIR                      = "config-dir"
	FLAG_EMAIL                           = "email"
	FLAG_REGISTER_UNSAFELY_WITHOUT_EMAIL = "register-unsafely-without-email"
	FLAG_STANDALONE                      = "standalone"
	FLAG_WEBROOT                         = "webroot"
	FLAG_AUTHENTICATOR                   = "authenticator"
	FLAG_AUTHENTICATOR_SHORT             = "a"
	FLAG_DOMAIN                          = "domain"
	FLAG_DOMAINS                         = "domains"
	FLAG_DOMAIN_SHORT                    = "d"
	FLAG_CERT_NAME                       = "cert-name"
)

var (
	flagHelp = &cli.Flag{
		Name:       FLAG_HELP,
		AltNames:   []string{FLAG_HELP_SHORT},
		TakesValue: true,
		HelpTopics: []string{TOPIC_OPTIONAL},
		Usage: cli.ArgumentUsage{
			Description: "show this help message and exit",
		},
		PostParse: func(f *cli.Flag, sc *cli.SubCommand, app *cli.App) error {
			if !f.IsPresent() {
				return nil
			}
			app.PrintHelp(f.StringOrDefault())
			return cli.ErrExitSuccess
		},
	}
	flagConfig = &cli.Flag{
		Name:         FLAG_CONFIG,
		AltNames:     []string{FLAG_CONFIG_SHORT},
		DefaultValue: &cli.SimpleValue{Value: defaultConfigFiles},
		TakesValue:   true,
		HelpTopics:   []string{TOPIC_OPTIONAL},
		Usage: cli.ArgumentUsage{
			ArgName:     "CONFIG_FILE",
			Description: "path to config file",
		},
	}
	flagWorkDir = &cli.Flag{
		Name:         FLAG_WORK_DIR,
		DefaultValue: &cli.SimpleValue{Value: defaultWorkDir},
		TakesValue:   true,
		HelpTopics:   []string{TOPIC_OPTIONAL},
	}
	flagLogsDir = &cli.Flag{
		Name:         FLAG_LOGS_DIR,
		DefaultValue: &cli.SimpleValue{Value: defaultLogsDir},
		TakesValue:   true,
		HelpTopics:   []string{TOPIC_OPTIONAL},
	}
	flagConfigDir = &cli.Flag{
		Name:         FLAG_CONFIG_DIR,
		DefaultValue: &cli.SimpleValue{Value: defaultConfigDir},
		TakesValue:   true,
		HelpTopics:   []string{TOPIC_OPTIONAL},
	}

	flagEmail = &cli.Flag{
		Name: FLAG_EMAIL,
		DefaultValue: &cli.AskValue{
			Query:  "Enter email address (used for urgent renewal and security notices)",
			Cancel: "An e-mail address or --register-unsafely-without-email must be provided.",
		},
		TakesValue: true,
	}
	flagRegisterUnsafely = &cli.Flag{
		Name:         FLAG_REGISTER_UNSAFELY_WITHOUT_EMAIL,
		DefaultValue: &cli.SimpleValue{Value: false},
		TakesValue:   true,
	}

	flagStandAlone = &cli.Flag{
		Name:         FLAG_STANDALONE,
		DefaultValue: &cli.SimpleValue{Value: false},
	}
	flagWebRoot = &cli.Flag{
		Name:         FLAG_WEBROOT,
		DefaultValue: &cli.SimpleValue{Value: false},
	}
	flagAuthenticator = &cli.Flag{
		Name:     FLAG_AUTHENTICATOR,
		AltNames: []string{FLAG_AUTHENTICATOR_SHORT},
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
	flagDomains = &cli.Flag{
		Name:     FLAG_DOMAIN,
		AltNames: []string{FLAG_DOMAINS, FLAG_DOMAIN_SHORT},
		DefaultValue: &cli.AskValue{
			Query:  "", // TODO
			Cancel: "", // TODO
		},
		TakesValue:    true,
		TakesMultiple: true,
		HelpTopics:    []string{COMMAND_CERTIFICATES},
		Usage: cli.ArgumentUsage{
			ArgName:     "DOMAIN",
			Description: "Domain names to apply. For multiple domains you can use multiple -d flags or enter a comma separated list of domains as a parameter. The first domain provided will be the subject CN of the certificate, and all domains will be Subject Alternative Names on the certificate. The first domain will also be used in some software user interfaces and as the file paths for the certificate and related material unless otherwise specified or you already have a certificate with the same name. In the case of a name collision it will append a number like 0001 to the file path name.",
		},
		PostParse: cli.RequireValueIfSet(),
	}
	flagCertName = &cli.Flag{
		Name: FLAG_CERT_NAME,
		DefaultValue: cli.SimpleValue{
			Value:        nil,
			UsageDefault: "the first provided domain or the name of an existing certificate on your system for the same domains",
		},
		TakesValue: true,
		HelpTopics: []string{COMMAND_CERTIFICATES},
		Usage: cli.ArgumentUsage{
			ArgName:     "CERTNAME",
			Description: "Certificate name to apply. This name is used by Certbot for housekeeping and in file paths; it doesn't affect the content of the certificate itself. To see certificate names, run 'certbot certificates'. When creating a new certificate, specifies the new certificate's name.",
		},
		PostParse: cli.RequireValueIfSet(),
	}
)
