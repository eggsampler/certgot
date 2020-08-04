package main

import (
	"fmt"
	"os"

	"github.com/eggsampler/certgot/util"

	"github.com/eggsampler/certgot/cli"
	"github.com/eggsampler/certgot/log"
)

func doPreRun(app *cli.App) error {
	ll := log.Entry{}
	for _, v := range app.GetArguments() {
		if !v.HasValue() {
			continue
		}
		ll = ll.WithField(v.Name, v.Value())
	}
	ll.Debug("cli arguments")
	if app.SpecificSubCommand != nil {
		log.WithField("subcommand", app.SpecificSubCommand.Name).Debug("cli subcommand")
	} else {
		log.Debug("no cli subcommand")
	}

	dirs := getDirectories()

	if err := setupDirectories(dirs); err != nil {
		return fmt.Errorf("Error setting up directories: %v\n", err)
	}

	if err := setupLocks(dirs); err != nil {
		return fmt.Errorf("Error setting up lock files: %v\n", err)
	}

	return nil
}

func getDirectories() []string {
	return []string{
		argConfigDir.StringOrDefault(),
		argLogsDir.StringOrDefault(),
		argWorkDir.StringOrDefault(),
	}
}

func setupDirectories(dirs []string) error {
	for _, dir := range dirs {
		if err := makeAndCheck(dir, 0755, os.Geteuid(), false); err != nil {
			return err
		}
	}
	return nil
}

func makeAndCheck(dir string, mode os.FileMode, uid int, strict bool) error {
	if err := os.MkdirAll(dir, mode); err != nil {
		return fmt.Errorf("error making directory %s - %v\n", dir, err)
	}
	if strict {
		fi, err := os.Stat(dir)
		if err != nil {
			return fmt.Errorf("error checking directory %s - %v\n", dir, err)
		}
		if fi.Mode() != mode {
			if err := os.Chmod(dir, mode); err != nil {
				return fmt.Errorf("error changing permission on directory %s to %d", dir, mode)
			}
		}
		if err := util.CheckUID(dir, fi.Sys(), uid); err != nil {
			return fmt.Errorf("error checking/changing owner on directory %s to %d", dir, uid)
		}
	}
	return nil
}
