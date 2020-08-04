package main

import (
	"github.com/eggsampler/certgot/cli"
	"github.com/eggsampler/certgot/log"
)

func doPostRun(_ *cli.App, r interface{}) {
	cleanupLocks()

	if r != nil {
		log.WithField("error", r).Debug("post run")
	} else {
		log.Debug("post run")
	}
}
