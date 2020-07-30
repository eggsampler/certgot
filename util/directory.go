package util

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func ParsePath(path string) string {
	if !strings.HasPrefix(path, "~") {
		return filepath.Clean(path)
	}
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome + path[1:])
	}
	if u, err := user.Current(); err != nil {
		return filepath.Join(u.HomeDir, path[1:])
	}
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, path[1:])
	}
	if home := filepath.Join(os.Getenv("HomeDrive"), os.Getenv("HomePath")); home != "" {
		return filepath.Join(home, path[1:])
	}
	if home := os.Getenv("UserProfile"); home != "" {
		return filepath.Join(home, path[1:])
	}
	return filepath.Clean(path)
}

func SetupDirectories(dirs []string) error {
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
		if err := checkUID(dir, fi.Sys(), uid); err != nil {
			return fmt.Errorf("error checking/changing owner on directory %s to %d", dir, uid)
		}
	}
	return nil
}
