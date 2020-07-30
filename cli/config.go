package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/eggsampler/certgot/util"
)

var (
	configLine = regexp.MustCompile(`([a-zA-Z\-]+)(?:\s*=\s*(.+))?`)
)

func (app *App) LoadConfig(configFiles []string) error {
	cfg := map[string]configEntry{}
	for _, fileName := range configFiles {
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
	fileName = util.ParsePath(fileName)
	if !util.FileExists(fileName) {
		return fmt.Errorf("config file doesn't exist: %s", fileName)
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
				return fmt.Errorf("error setting argument %s to value %q: %v", key, entry.value, err)
			}
		}
	}
	return nil
}
