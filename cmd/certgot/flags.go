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
	FLAG_NON_INTERACTIVE                 = "non-interactive"
	FLAG_NONINTERACTIVE                  = "noninteractive"
	FLAG_NON_INTERACTIVE_SHORT           = "n"
	FLAG_FORCE_INTERACTIVE               = "force-interactive"
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
			topic, _ := f.String(getFlagValues(true))
			app.PrintHelp(topic)
			return cli.ErrExitSuccess
		},
	}
	flagConfig = &cli.Flag{
		Name:       FLAG_CONFIG,
		AltNames:   []string{FLAG_CONFIG_SHORT},
		Value:      cli.NewSimpleValueDefault(defaultConfigFiles),
		TakesValue: true,
		HelpTopics: []string{TOPIC_OPTIONAL},
		Usage: cli.ArgumentUsage{
			ArgName:     "CONFIG_FILE",
			Description: "path to config file",
		},
	}
	flagWorkDir = &cli.Flag{
		Name:       FLAG_WORK_DIR,
		Value:      cli.NewSimpleValueDefault(defaultWorkDir),
		TakesValue: true,
		HelpTopics: []string{TOPIC_OPTIONAL},
	}
	flagLogsDir = &cli.Flag{
		Name:       FLAG_LOGS_DIR,
		Value:      cli.NewSimpleValueDefault(defaultLogsDir),
		TakesValue: true,
		HelpTopics: []string{TOPIC_OPTIONAL},
	}
	flagConfigDir = &cli.Flag{
		Name:       FLAG_CONFIG_DIR,
		Value:      cli.NewSimpleValueDefault(defaultConfigDir),
		TakesValue: true,
		HelpTopics: []string{TOPIC_OPTIONAL},
	}

	flagEmail = &cli.Flag{
		Name: FLAG_EMAIL,
		Value: &cli.AskValue{
			Query:  "Enter email address (used for urgent renewal and security notices)",
			Cancel: "An e-mail address or --register-unsafely-without-email must be provided.",
		},
		TakesValue: true,
	}
	flagRegisterUnsafely = &cli.Flag{
		Name:       FLAG_REGISTER_UNSAFELY_WITHOUT_EMAIL,
		Value:      cli.NewSimpleValueDefault(false),
		TakesValue: true,
	}

	flagStandAlone = &cli.Flag{
		Name:  FLAG_STANDALONE,
		Value: cli.NewSimpleValueDefault(false),
	}
	flagWebRoot = &cli.Flag{
		Name:  FLAG_WEBROOT,
		Value: cli.NewSimpleValueDefault(false),
	}
	flagAuthenticator = &cli.Flag{
		Name:     FLAG_AUTHENTICATOR,
		AltNames: []string{FLAG_AUTHENTICATOR_SHORT},
		Value: &cli.ListValue{
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
		Value: &cli.AskValue{
			Query:  "", // TODO
			Cancel: "", // TODO
		},
		TakesValue:    true,
		TakesMultiple: true,
		HelpTopics:    []string{CMD_CERTIFICATES},
		Usage: cli.ArgumentUsage{
			ArgName:     "DOMAIN",
			Description: "Domain names to apply. For multiple domains you can use multiple -d flags or enter a comma separated list of domains as a parameter. The first domain provided will be the subject CN of the certificate, and all domains will be Subject Alternative Names on the certificate. The first domain will also be used in some software user interfaces and as the file paths for the certificate and related material unless otherwise specified or you already have a certificate with the same name. In the case of a name collision it will append a number like 0001 to the file path name.",
		},
		PostParse: cli.RequireValueIfSet(),
	}
	flagCertName = &cli.Flag{
		Name:       FLAG_CERT_NAME,
		Value:      cli.NewSimpleValueNotSet("the first provided domain or the name of an existing certificate on your system for the same domains"),
		TakesValue: true,
		HelpTopics: []string{CMD_CERTIFICATES},
		Usage: cli.ArgumentUsage{
			ArgName:     "CERTNAME",
			Description: "Certificate name to apply. This name is used by Certbot for housekeeping and in file paths; it doesn't affect the content of the certificate itself. To see certificate names, run 'certbot certificates'. When creating a new certificate, specifies the new certificate's name.",
		},
		PostParse: cli.RequireValueIfSet(),
	}
	flagNonInteractive = &cli.Flag{
		Name:     FLAG_NON_INTERACTIVE,
		AltNames: []string{FLAG_NONINTERACTIVE, FLAG_NON_INTERACTIVE_SHORT},
		Value:    cli.NewSimpleValueDefault(false),
		Usage: cli.ArgumentUsage{
			Description: "Run without ever asking for user input. This may require additional command line flags; the client will try to explain which ones are required if it finds one missing",
		},
	}
	flagForceInteractive = &cli.Flag{
		Name:       FLAG_FORCE_INTERACTIVE,
		Value:      cli.NewSimpleValueDefault(false),
		HelpTopics: []string{CMD_CERTONLY},
		Usage: cli.ArgumentUsage{
			Description: "Force Certbot to be interactive even if it detects it's not being run in a terminal. This flag cannot be used with the renew subcommand.",
		},
	}
)

// TODO: not super happy how this turned out, investigate some way to separate values from flags?
// maybe separate the value on a flag from "config values"
func getFlagValues(incDef bool) (nonInteractive, forceInteractive, isTerminal, includeDefault bool) {
	nonInteractive, _ = flagNonInteractive.Bool(true, false, false, false)
	forceInteractive, _ = flagForceInteractive.Bool(true, false, false, false)
	isTerminal = cli.IsTerminal()
	includeDefault = incDef
	return
}
