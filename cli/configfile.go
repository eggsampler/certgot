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
	// TODO: ini reader?
	configLine = regexp.MustCompile(`^([a-zA-Z\-]+)(?:\s*=\s*(.+))?$`)
)

type configFileEntry struct {
	fileName string
	line     int
	key      string
	hasValue bool
	value    string
}

func parseConfig(r io.Reader, fileName string) (map[string][]configFileEntry, error) {
	var entries map[string][]configFileEntry
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
		entry := configFileEntry{
			fileName: fileName,
			line:     lineNumber,
			key:      m[1],
		}
		if strings.Contains(m[0], "=") {
			entry.hasValue = true
			entry.value = m[2]
		}
		if entries == nil {
			entries = map[string][]configFileEntry{}
		}
		entries[entry.key] = append(entries[entry.key], entry)
	}
	return entries, scanner.Err()
}

func setConfig(entries map[string][]configFileEntry, cl ConfigList) error {
	for name, entryList := range entries {
		cfg := cl.Get(name)
		if cfg == nil {
			return fmt.Errorf("unknown config %q on line %d in config file: %s", name, entryList[0].line, entryList[0].fileName)
		}

		var values []string

		for _, v := range entryList {
			ll := log.WithField("filename", v.fileName).
				WithField("line", v.line).
				WithField("config", v.key).
				WithField("hasValue", v.hasValue)
			ll.Trace("config present")

			values = append(values, v.value)
		}

		log.WithField("config", name).WithField("value", entryList).Trace("setting value")
		if err := cfg.set(values); err != nil {
			return fmt.Errorf("error setting config %q to value %q: %v", name, values, err)
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

func loadConfig(configFiles []string, skipOpenErrors bool, cl ConfigList, sys fs.FS) error {
	if len(configFiles) == 0 {
		return errors.New("no config files provided")
	}

	entries := map[string][]configFileEntry{}

	for _, fileName := range configFiles {
		fileName = parsePath(fileName, os.Getenv, user.Current)
		ll := log.WithField("filename", fileName)
		ll.Trace("attempting to read config file")
		// TODO: why does io/fs.ValidPath return false for paths starting/ending with a slash???
		f, err := sys.Open(strings.Trim(fileName, string(os.PathSeparator)))
		if err != nil {
			// skip file errors if config file isn't explicitly set
			if skipOpenErrors {
				ll.WithError(err).Debug("reading config file")
				continue
			}
			return fmt.Errorf("error opening config file %q: %w", fileName, err)
		}
		ll.Trace("parsing config file")
		if c, err := parseConfig(f, fileName); err != nil {
			return err
		} else {
			for k, v := range c {
				entries[k] = append(entries[k], v...)
			}
		}
		_ = f.Close()
	}
	return setConfig(entries, cl)
}
