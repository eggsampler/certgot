package log

import (
	"fmt"
)

// why does every structured logging library need like 500 non-vendored deps?!?!
// TODO: replace this with something good or fix it up

type Level int

const (
	DEBUG Level = iota
	INFO
	ERROR
	FATAL
)

func SetLevel(level Level) {
	currentLevel = level
}

func GetLevel() Level {
	return currentLevel
}

func SetOutputFile(file string) {
	logFile = file
}

func WithError(err error) Entry {
	return Entry{parseKV("error", err)}
}

func WithField(key string, value interface{}) Entry {
	return Entry{parseKV(key, value)}
}

func WithFields(fields ...interface{}) Entry {
	entry := make(Entry, len(fields)/2)
	if len(entry) == 0 {
		return entry
	}
	for i := 0; i < len(fields); i += 2 {
		if i+1 == len(fields) {
			break
		}
		entry = append(entry, parseKV(fields[i], fields[i+1]))
	}
	return entry
}

func Printf(format string, a ...interface{}) {
	logInfo(fmt.Sprintf(format, a...), nil)
}

func Log(level Level, message string) {
	log(level, message, nil)
}

func Debug(message string) {
	logDebug(message, nil)
}

func Info(message string) {
	logInfo(message, nil)
}

func Error(message string) {
	logError(message, nil)
}

func Fatal(message string) {
	logFatal(message, nil)
}
