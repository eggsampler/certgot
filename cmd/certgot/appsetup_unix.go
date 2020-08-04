package main

import "path/filepath"

var (
	defaultConfigFiles = []string{
		string(filepath.Separator) + filepath.Join("etc", "letsencrypt", "cli.ini"),
		filepath.Join("~", ".config", "letsencrypt", "cli.ini"),
	}
	defaultWorkDir   = filepath.Join(string(filepath.Separator), "var", "lib", "letsencrypt")
	defaultLogsDir   = filepath.Join(string(filepath.Separator), "var", "logs", "letsencrypt")
	defaultConfigDir = filepath.Join(string(filepath.Separator), "etc", "letsencrypt")
)
