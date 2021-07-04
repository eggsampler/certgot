// +build windows

package main

import "path/filepath"

var (
	// TODO: check what the cli.parsePath function evalutes ~ to on windows
	defaultConfigFiles = []string{
		filepath.Join("c:", "certbot", "cli.ini"),
		filepath.Join("~", "certbot", "cli.ini"),
	}
	defaultWorkDir   = filepath.Join("c:", "certbot", "lib")
	defaultLogsDir   = filepath.Join("c:", "certbot", "log")
	defaultConfigDir = filepath.Join("c:", "certbot")
)
