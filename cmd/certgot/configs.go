package main

import (
	"fmt"

	"github.com/eggsampler/certgot/cli"
)

const (
	CONFIG_FILE = "filename"
)

var (
	cfgConfigFile = &cli.Config{
		Name:        CONFIG_FILE,
		Default:     defaultConfigFiles,
		HelpDefault: fmt.Sprintf("%+v", defaultConfigFiles),
	}
)
