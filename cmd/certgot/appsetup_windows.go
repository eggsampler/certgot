package main

import "path/filepath"

var (
	defaultConfigFiles = []string{
		filepath.Join("c:", "certbot", "cli.ini"),
	}
	defaultWorkDir   = filepath.Join("c:", "certbot", "lib")
	defaultLogsDir   = filepath.Join("c:", "certbot", "log")
	defaultConfigDir = filepath.Join("c:", "certbot")
)
