package main

import "github.com/eggsampler/certgot/cli"

// TODO: pick a better naming scheme to identify the constant names vs the variable arguments
// to better match go naming https://golang.org/doc/effective_go#mixed-caps

const (
	ARG_HELP                            = "help"
	ARG_HELP_SHORT                      = "h"
	ARG_CONFIG                          = "config"
	ARG_CONFIG_SHORT                    = "c"
	ARG_WORK_DIR                        = "work-dir"
	ARG_LOGS_DIR                        = "logs-dir"
	ARG_CONFIG_DIR                      = "config-dir"
	ARG_EMAIL                           = "email"
	ARG_REGISTER_UNSAFELY_WITHOUT_EMAIL = "register-unsafely-without-email"
	ARG_STANDALONE                      = "standalone"
	ARG_WEBROOT                         = "webroot"
	ARG_AUTHENTICATOR                   = "authenticator"
	ARG_AUTHENTICATOR_SHORT             = "a"
	ARG_DOMAIN                          = "domain"
	ARG_DOMAINS                         = "domains"
	ARG_DOMAIN_SHORT                    = "d"
	ARG_CERT_NAME                       = "cert-name"
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
			if !arg.IsPresent() {
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
		Name:         ARG_WORK_DIR,
		DefaultValue: &cli.SimpleValue{Value: defaultWorkDir},
		TakesValue:   true,
		HelpTopics:   []string{TOPIC_OPTIONAL},
	}
	argLogsDir = &cli.Argument{
		Name:         ARG_LOGS_DIR,
		DefaultValue: &cli.SimpleValue{Value: defaultLogsDir},
		TakesValue:   true,
		HelpTopics:   []string{TOPIC_OPTIONAL},
	}
	argConfigDir = &cli.Argument{
		Name:         ARG_CONFIG_DIR,
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
	argDomains = &cli.Argument{
		Name:     ARG_DOMAIN,
		AltNames: []string{ARG_DOMAINS, ARG_DOMAIN_SHORT},
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
	argCertName = &cli.Argument{
		Name: ARG_CERT_NAME,
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
