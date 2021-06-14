package main

import (
	"fmt"
	"runtime"
)

func init() {
	if len(defaultConfigFiles) == 0 {
		panic(fmt.Sprintf("no default config files set for os %q arch %q", runtime.GOOS, runtime.GOARCH))
	}
	if defaultWorkDir == "" {
		panic(fmt.Sprintf("no default work dir set for os %q arch %q", runtime.GOOS, runtime.GOARCH))
	}
	if defaultLogsDir == "" {
		panic(fmt.Sprintf("no default logs dir set for os %q arch %q", runtime.GOOS, runtime.GOARCH))
	}
	if defaultConfigDir == "" {
		panic(fmt.Sprintf("no default config dir set for os %q arch %q", runtime.GOOS, runtime.GOARCH))
	}
}
