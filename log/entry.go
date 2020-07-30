package log

import (
	"fmt"
)

type Entry []string

func (l Entry) WithError(err error) Entry {
	return append(l, fmt.Sprintf("error=%q", err))
}

func (l Entry) WithField(key string, value interface{}) Entry {
	return append(l, parseKV(key, value))
}

func (l Entry) WithFields(fields ...interface{}) Entry {
	return append(l, WithFields(fields)...)
}

func (l Entry) Printf(format string, a ...interface{}) {
	logInfo(fmt.Sprintf(format, a...), l)
}

func (l Entry) Log(level Level, msg string) {
	log(level, msg, l)
}

func (l Entry) Debug(msg string) {
	logDebug(msg, l)
}

func (l Entry) Info(msg string) {
	logInfo(msg, l)
}

func (l Entry) Error(msg string) {
	logError(msg, l)
}

func (l Entry) Fatal(msg string) {
	logFatal(msg, l)
}
