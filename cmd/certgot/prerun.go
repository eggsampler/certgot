package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/eggsampler/certgot/cli"
	"github.com/eggsampler/certgot/log"
	"github.com/eggsampler/certgot/util"
)

func doPreRun(ctx *cli.Context) error {
	cfg := ctx.App.Configs.Get(CONFIG_FILE)
	if err := ctx.App.LoadConfig(cfg.StringSlice(), !cfg.IsSet()); err != nil {
		return err
	}

	ll := log.Fields{}
	for _, v := range ctx.Flags {
		if len(v.StringSlice()) > 0 {
			ll = ll.WithField(v.Name, v.StringSlice())
		} else {
			ll = ll.WithField(v.Name, "(no value)")
		}
	}
	ll.Debug("cli arguments")
	if ctx.Command != nil {
		log.WithField("command", ctx.Command.Name).Debug("cli command")
	} else {
		log.Debug("no cli command")
	}

	dirs, err := getDirectories(ctx)
	if err != nil {
		return fmt.Errorf("error fetching directories: %w", err)
	}

	for _, dir := range dirs {
		if err := makeAndCheck(dir, 0755, os.Geteuid(), false); err != nil {
			return err
		}
	}

	if err := setupLocks(dirs); err != nil {
		return fmt.Errorf("error setting up lock files: %w", err)
	}

	return nil
}

func getDirectories(ctx *cli.Context) ([]string, error) {
	cd := cfgConfigDir.String()
	if len(cd) == 0 {
		return nil, errors.New("no config directory set")
	}
	ld := cfgLogsDir.String()
	if len(cd) == 0 {
		return nil, errors.New("no logs directory set")
	}
	wd := cfgWorkDir.String()
	if len(cd) == 0 {
		return nil, errors.New("no work directory set")
	}
	return []string{
		cd,
		ld,
		wd,
	}, nil
}

func makeAndCheck(dir string, mode os.FileMode, uid int, strict bool) error {
	log.WithFields("dir", dir, "mode", mode, "uid", uid, "strict", strict).Trace("making directory")

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
