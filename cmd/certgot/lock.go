package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/eggsampler/certgot/log"

	"github.com/gofrs/flock"
)

var (
	lockFiles []*flock.Flock
)

func setupLocks(dirs []string) error {
	for _, dir := range dirs {
		lf := flock.New(filepath.Join(dir, ".certbot.lock"))
		log.WithField("lockfile", lf.Path()).Trace("creating lock file")
		locked, err := lf.TryLock()
		if err != nil {
			return fmt.Errorf("trying lock file %s - %v\n", lf.Path(), err)
		}
		if !locked {
			return fmt.Errorf("obtaining lock file %s - %v\n", lf.Path(), err)
		}
		lockFiles = append(lockFiles, lf)
	}
	return nil
}

func cleanupLocks() (errors []error) {
	for _, lf := range lockFiles {
		log.WithField("lockfile", lf.Path()).Trace("cleaning up lock file")
		if err := lf.Unlock(); err != nil {
			errors = append(errors, fmt.Errorf("error unlocking file %q: %w", lf.Path(), err))
		}
		if err := os.Remove(lf.Path()); err != nil {
			errors = append(errors, fmt.Errorf("error removing file %q: %w", lf.Path(), err))
		}
	}

	return
}
