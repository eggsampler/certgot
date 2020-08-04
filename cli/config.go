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
	configLine = regexp.MustCompile(`([a-zA-Z\-]+)(?:\s*=\s*(.+))?`)
)

func (app *App) LoadConfig(argCfg *Argument) error {
	cfg := map[string]configEntry{}
	for _, fileName := range argCfg.StringSliceOrDefault() {
		// skip file not found errors if config is default cfg files
		if !argCfg.isPresent && !fileExists(fileName) {
			continue
		}
		if err := parseConfigFile(cfg, fileName); err != nil {
			return err
		}
	}
	return setConfig(cfg, app.args)
}

type configEntry struct {
	fileName string
	line     int
	key      string
	hasValue bool
	value    string
}

func parseConfigFile(cfg map[string]configEntry, fileName string) error {
	fileName = parsePath(fileName)
	if !fileExists(fileName) {
		return fmt.Errorf("error loading config file %s: %w", fileName, os.ErrNotExist)
	}
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening config file %s - %v", fileName, err)
	}
	defer f.Close()
	return parseConfig(cfg, f, fileName)
}

func parseConfig(cfg map[string]configEntry, r io.Reader, fileName string) error {
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
			return fmt.Errorf("invalid argument %q on line %d in config file: %s", line, lineNumber, fileName)
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
		cfg[entry.key] = entry
	}
	return nil
}

func setConfig(config map[string]configEntry, args map[string]*Argument) error {
	for key, entry := range config {
		arg, ok := args[key]
		if !ok {
			return fmt.Errorf("unknown argument %s on line %d in config file: %s", key, entry.line, entry.fileName)
		}
		arg.isPresent = true
		if entry.hasValue {
			if err := arg.Set(entry.value); err != nil {
				return fmt.Errorf("error setting argument %q to value %q: %v", key, entry.value, err)
			}
		}
	}
	return nil
}

func parsePath(path string) string {
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

func fileExists(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err) || !os.IsNotExist(err)
}
