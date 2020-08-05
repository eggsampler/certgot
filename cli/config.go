package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
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
	return cfg, nil
}

func setConfig(config []configEntry, args map[string]*Argument) error {
	for _, entry := range config {
		arg, ok := args[entry.key]
		if !ok {
			return fmt.Errorf("unknown argument %s on line %d in config file: %s", entry.key, entry.line, entry.fileName)
		}
		arg.isPresent = true
		if entry.hasValue {
			if err := arg.Set(entry.value); err != nil {
				return fmt.Errorf("error setting argument %q to value %q: %v", entry.key, entry.value, err)
			}
		}
	}
	return nil
}

// parsePath takes a path string which might begin with a ~ and attempts to replace it with the users home directory
func parsePath(path string) string {
	if !strings.HasPrefix(path, "~") {
		return filepath.Clean(path)
	}
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome + path[1:])
	}
	if u, err := user.Current(); err == nil {
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
