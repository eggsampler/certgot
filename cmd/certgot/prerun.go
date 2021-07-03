package main

import (
	"github.com/eggsampler/certgot/cli"
	"github.com/eggsampler/certgot/log"
)

func doPreRun(ctx *cli.Context) error {
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
		log.WithField("command", ctx.Command.Name).Debug("cli subcommand")
	} else {
		log.Debug("no cli subcommand")
	}

	/*

		dirs, err := getDirectories()
		if err != nil {
			return fmt.Errorf("Error fetching directories: %w", err)
		}

		if err := setupDirectories(dirs); err != nil {
			return fmt.Errorf("Error setting up directories: %v", err)
		}

		if err := setupLocks(dirs); err != nil {
			return fmt.Errorf("Error setting up lock files: %v", err)
		}

	*/

	return nil
}

/*

func getDirectories() ([]string, error) {
	cd, err := flagConfigDir.String(true, false, false, true)
	if err != nil {
		return nil, err
	}
	ld, err := flagLogsDir.String(true, false, false, true)
	if err != nil {
		return nil, err
	}
	wd, err := flagWorkDir.String(true, false, false, true)
	if err != nil {
		return nil, err
	}
	return []string{
		cd,
		ld,
		wd,
	}, nil
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


*/
