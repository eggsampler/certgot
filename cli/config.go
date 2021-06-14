package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/eggsampler/certgot/log"
)

var (
	configLine = regexp.MustCompile(`^([a-zA-Z\-]+)(?:\s*=\s*(.+))?$`)
)

type configEntry struct {
	fileName string
	line     int
	key      string
	hasValue bool
	value    string
}

func parseConfig(r io.Reader, fileName string) ([]configEntry, error) {
	var cfg []configEntry
	lineNumber := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		m := configLine.FindStringSubmatch(line)
		if m == nil {
			return nil, fmt.Errorf("invalid argument %q on line %d in config file: %s", line, lineNumber, fileName)
		}
		entry := configEntry{
			fileName: fileName,
			line:     lineNumber,
			key:      m[1],
		}
		if strings.Contains(m[0], "=") {
			entry.hasValue = true
			entry.value = m[2]
		}
		cfg = append(cfg, entry)
	}
	return cfg, scanner.Err()
}

func setConfig(config []configEntry, args map[string]*Argument) error {
	for _, entry := range config {
		arg, ok := args[entry.key]
		if !ok {
			return fmt.Errorf("unknown argument %s on line %d in config file: %s", entry.key, entry.line, entry.fileName)
		}
		ll := log.WithField("filename", entry.fileName).
			WithField("line", entry.line).
			WithField("arg", entry.key).
			WithField("hasValue", entry.hasValue)
		ll.Trace("arg present")
		arg.IsPresent = true
		if entry.hasValue {
			ll.WithField("value", entry.value).Trace("setting value")
			if err := arg.Set(entry.value); err != nil {
				return fmt.Errorf("error setting argument %q to value %q: %v", entry.key, entry.value, err)
			}
		}
	}
	return nil
}

// parsePath takes a path string which might begin with a ~ and, if it does, attempts to replace the tilde
// with the users home directory
// Priorities:
//  $XDG_CONFIG_HOME
//  user.Current().HomeDir
//  $HOME
//  %HomeDrive%\%HomePath% (windows)
//  %UserProfile% (windows)
func parsePath(path string, envFunc func(string) string, userFunc func() (*user.User, error)) string {
	if !strings.HasPrefix(path, "~") {
		return filepath.Clean(path)
	}
	if xdgConfigHome := envFunc("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome + path[1:])
	}
	if u, err := userFunc(); err == nil {
		return filepath.Join(u.HomeDir, path[1:])
	}
	if home := envFunc("HOME"); home != "" {
		return filepath.Join(home, path[1:])
	}
	if home := filepath.Join(envFunc("HomeDrive"), envFunc("HomePath")); home != "" {
		return filepath.Join(home, path[1:])
	}
	if home := envFunc("UserProfile"); home != "" {
		return filepath.Join(home, path[1:])
	}
	return filepath.Clean(path)
}

func loadConfig(app *App, cfgFile *Argument, sys fs.FS) error {
	if cfgFile == nil {
		return errors.New("no config file argument provided")
	}
	var cfg []configEntry
	for _, fileName := range cfgFile.StringSliceOrDefault() {
		fileName = parsePath(fileName, os.Getenv, user.Current)
		ll := log.WithField("filename", fileName)
		ll.Trace("attempting to read config file")
		// TODO: why does io/fs.ValidPath return false for paths starting/ending with a slash???
		f, err := sys.Open(strings.Trim(fileName, string(os.PathSeparator)))
		if err != nil {
			ll.WithError(err).Error("reading config file")
			// skip file errors if config file isn't explicitly set
			if !cfgFile.IsPresent {
				continue
			}
			return fmt.Errorf("error opening config file %q: %w", fileName, err)
		}
		ll.Trace("parsing config file")
		if c, err := parseConfig(f, fileName); err != nil {
			return err
		} else {
			cfg = append(cfg, c...)
		}
		_ = f.Close()
	}
	return setConfig(cfg, app.argsMap)
}
