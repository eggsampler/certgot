package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
)

var (
	lockFiles []*flock.Flock
)

func setupLocks(dirs []string) error {
	for _, dir := range dirs {
		lf := flock.New(filepath.Join(dir, ".certbot.lock"))
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

func cleanupLocks() {
	for _, lf := range lockFiles {
		_ = lf.Unlock()
		_ = os.Remove(lf.Path())
	}
}
