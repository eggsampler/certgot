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
		Name:            FLAG_HELP,
		AltNames:        []string{FLAG_HELP_SHORT},
		TakesValue:      true,
		HelpCategories:  []string{CATEGORY_OPTIONAL},
		HelpDescription: "show this help message and exit",
		PostParseFunc: func(f *cli.Flag, ctx *cli.Context) error {
			_ = ctx.App.PrintHelp(ctx, f.String())
			return cli.ErrExitSuccess
		},
	}
	flagConfigFile = &cli.Flag{
		Name:            FLAG_CONFIG,
		AltNames:        []string{FLAG_CONFIG_SHORT},
		TakesValue:      true,
		RequiresValue:   true,
		HelpDefault:     cli.GetConfigDefault(CONFIG_FILE),
		HelpCategories:  []string{CATEGORY_OPTIONAL},
		HelpValueName:   "CONFIG_FILE",
		HelpDescription: "path to config file",
		PostParseFunc:   cli.SetConfigValue(CONFIG_FILE),
	}
	flagWorkDir = &cli.Flag{
		Name:           FLAG_WORK_DIR,
		TakesValue:     true,
		RequiresValue:  true,
		HelpCategories: []string{CATEGORY_OPTIONAL},
	}
	flagLogsDir = &cli.Flag{
		Name:           FLAG_LOGS_DIR,
		TakesValue:     true,
		RequiresValue:  true,
		HelpCategories: []string{CATEGORY_OPTIONAL},
	}
	flagConfigDir = &cli.Flag{
		Name:           FLAG_CONFIG_DIR,
		TakesValue:     true,
		RequiresValue:  true,
		HelpCategories: []string{CATEGORY_OPTIONAL},
	}

	flagDomains = &cli.Flag{
		Name:            FLAG_DOMAIN,
		AltNames:        []string{FLAG_DOMAINS, FLAG_DOMAIN_SHORT},
		TakesValue:      true,
		RequiresValue:   true,
		AllowMultiple:   true,
		HelpCategories:  []string{CMD_CERTIFICATES},
		HelpValueName:   "DOMAIN",
		HelpDescription: "Domain names to apply. For multiple domains you can use multiple -d flags or enter a comma separated list of domains as a parameter. The first domain provided will be the subject CN of the certificate, and all domains will be Subject Alternative Names on the certificate. The first domain will also be used in some software user interfaces and as the file paths for the certificate and related material unless otherwise specified or you already have a certificate with the same name. In the case of a name collision it will append a number like 0001 to the file path name.",
	}
	flagCertName = &cli.Flag{
		Name:           FLAG_CERT_NAME,
		TakesValue:     true,
		RequiresValue:  true,
		HelpCategories: []string{CMD_CERTIFICATES},
		HelpDefault: func(*cli.Context) (string, error) {
			return "the first provided domain or the name of an existing certificate on your system for the same domains", nil
		},
		HelpValueName:   "CERTNAME",
		HelpDescription: "Certificate name to apply. This name is used by Certbot for housekeeping and in file paths; it doesn't affect the content of the certificate itself. To see certificate names, run 'certbot certificates'. When creating a new certificate, specifies the new certificate's name.",
	}
	flagNonInteractive = &cli.Flag{
		Name:            FLAG_NON_INTERACTIVE,
		AltNames:        []string{FLAG_NONINTERACTIVE, FLAG_NON_INTERACTIVE_SHORT},
		HelpDescription: "Run without ever asking for user input. This may require additional command line flags; the client will try to explain which ones are required if it finds one missing",
	}
	flagForceInteractive = &cli.Flag{
		Name:            FLAG_FORCE_INTERACTIVE,
		HelpCategories:  []string{CMD_CERTONLY},
		HelpDescription: "Force Certbot to be interactive even if it detects it's not being run in a terminal. This flag cannot be used with the renew subcommand.",
	}
)
