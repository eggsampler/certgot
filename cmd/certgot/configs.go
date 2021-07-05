package main

import (
	"fmt"

	"github.com/eggsampler/certgot/cli"
)

const (
	CONFIG_FILE       = "filename"
	CONFIG_LOGS_DIR   = "logs-dir"
	CONFIG_CONFIG_DIR = "config-dir"
	CONFIG_WORK_DIR   = "work-dir"
	CONFIG_DOMAINS    = "domains"
	CONFIG_CERT_NAME  = "cert-name"
)

var (
	cfgConfigFile = &cli.Config{
		Name:        CONFIG_FILE,
		Default:     defaultConfigFiles,
		HelpDefault: fmt.Sprintf("%+v", defaultConfigFiles),
	}
	cfgLogsDir = &cli.Config{
		Name:        CONFIG_LOGS_DIR,
		Default:     []string{defaultLogsDir},
		HelpDefault: defaultLogsDir,
	}
	cfgConfigDir = &cli.Config{
		Name:        CONFIG_CONFIG_DIR,
		Default:     []string{defaultConfigDir},
		HelpDefault: defaultConfigDir,
	}
	cfgWorkDir = &cli.Config{
		Name:        CONFIG_WORK_DIR,
		Default:     []string{defaultWorkDir},
		HelpDefault: defaultWorkDir,
	}
	cfgDomains = &cli.Config{
		Name:        CONFIG_DOMAINS,
		Default:     nil,
		HelpDefault: "",
		OnSet:       nil,
	}
	cfgCertName = &cli.Config{
		Name:        CONFIG_CERT_NAME,
		Default:     nil,
		HelpDefault: "",
		OnSet:       nil,
	}
)
