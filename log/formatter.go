package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Formatter interface {
	FormatLog(level Level, msg string, fields []logField)
}

func DefaultFormatter() Formatter {
	return StringFormatter(os.Stdout)
}

var CurrentFormatter = DefaultFormatter()

func formatLog(level Level, msg string, fields []logField) {
	if level < currentLevel {
		return
	}

	CurrentFormatter.FormatLog(level, msg, fields)
}

type stringFormatter struct {
	w io.Writer
}

var stringFormatterOnce sync.Once
var stringFormatterInstance *stringFormatter

func StringFormatter(w io.Writer) Formatter {
	if stringFormatterInstance == nil {
		stringFormatterOnce.Do(
			func() {
				stringFormatterInstance = &stringFormatter{w}
			})
	}
	return stringFormatterInstance
}

func (sf stringFormatter) FormatLog(level Level, msg string, fields []logField) {
	_, _ = fmt.Fprintln(sf.w, stringFormatLogLine(level, msg, fields))
}

func stringFormatLogLine(level Level, msg string, fields []logField) string {
	line := fmt.Sprintf("%s[%s] {%s} msg=%q",
		levelNames[level], getTime(), getSource(), msg)
	if len(fields) > 0 {
		var s []string
		for _, f := range fields {
			s = append(s, f.Key()+"="+fmt.Sprintf("%q", f.Value()))
		}
		line = line + " " + strings.Join(s, " ")
	}
	return line
}
