package log

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	levels = map[Level]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		ERROR: "ERROR",
		FATAL: "FATAL",
	}

	currentLevel = Level(0)

	logFile       string
	logFileHandle *os.File
)

func init() {
	v, ok := os.LookupEnv("CERTBOT_LOGLEVEL")
	if !ok {
		return
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return
	}
	SetLevel(Level(i))
}

func parseKV(key, value interface{}, verbosity int) string {
	if verbosity <= 0 {
		return fmt.Sprintf(`%s="%v"`, key, value)
	} else if verbosity == 1 {
		return fmt.Sprintf(`%s="%+v"`, key, value)
	} // else if verbosity >= 2 {
	return fmt.Sprintf(`%s="%#v"`, key, value)
}

func log(level Level, msg string, fields []string) {
	fmsg := fmt.Sprintf("%s[%s] msg=%q %s",
		levels[level],
		time.Now().Format("2006-01-02 15:04:05 -0700"),
		msg,
		strings.Join(fields, " "))

	if level >= currentLevel {
		fmt.Println(fmsg)
	}

	if logFileHandle == nil && logFile != "" {
		var err error
		logFileHandle, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic(fmt.Sprintf("Error opening log file at %q: %v", logFile, err))
		}
	}

	if logFileHandle != nil {
		if _, err := fmt.Fprintln(logFileHandle, fmsg); err != nil {
			panic("Error writing to log: " + err.Error())
		}
	}
}

func logDebug(msg string, fields []string) {
	log(DEBUG, msg, fields)
}

func logInfo(msg string, fields []string) {
	log(INFO, msg, fields)
}

func logError(msg string, fields []string) {
	log(ERROR, msg, fields)
}

func logFatal(msg string, fields []string) {
	log(FATAL, msg, fields)
}
