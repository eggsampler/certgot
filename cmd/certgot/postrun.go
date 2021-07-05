package main

import (
	"github.com/eggsampler/certgot/cli"
	"github.com/eggsampler/certgot/log"
)

func doPostRun(ctx *cli.Context, err error) error {
	errs := cleanupLocks()
	for _, v := range errs {
		log.WithError(v).Debug("cleaning up locks")
	}

	if err != nil {
		log.WithField("error", err).Debug("post run")
	} else {
		log.Debug("post run")
	}

	return nil
}
