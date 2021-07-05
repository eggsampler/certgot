package log

import (
	"fmt"
	"os"
)

// Fields represents a list of pre-formatted log fields
// This type is designed to continue the method chaining which starts in log.go
// and also hold the logic for those functions
// You won't use this type directly, only instantiated from log.go functions
type Fields []logField

func (f Fields) WithError(err error) Fields {
	if err == nil {
		return f
	}
	newFields := append(f, errorLogField{err})
	if errorFields := getErrorFields(err); len(errorFields) > 0 {
		newFields = append(newFields, errorFields...)
	}
	return newFields
}

func (f Fields) WithField(key, value interface{}) Fields {
	return append(f, toField(key, value))
}

func toString(val interface{}, format string) string {
	switch v := val.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf(format, v)
	}
}

func toField(key interface{}, val interface{}) logField {
	var sVal string

	// TODO: is there any value in this changing?
	// because it changes based on the current log level
	// not the final .Debug vs .Trace vs .Info level
	lvl := int(GetLevel())
	switch {
	case lvl <= 0:
		sVal = toString(val, "%#v")
	case lvl == 1:
		sVal = toString(val, "%+v")
	default: // case lvl >= 2:
		sVal = toString(val, "%v")
	}

	return stringLogField{toString(key, "%v"), sVal}
}

func (f Fields) WithFields(fields ...interface{}) Fields {
	if len(fields) == 0 {
		return f
	}
	var ff []logField
	for i := 0; i < len(fields); i += 2 {
		ff = append(ff, toField(fields[i], fields[i+1]))
	}
	return append(f, ff...)
}

func (f Fields) Printf(format string, a ...interface{}) {
	formatLog(InfoLevel, fmt.Sprintf(format, a...), f)
}

func (f Fields) Log(level Level, msg string) {
	formatLog(level, msg, f)
}

func (f Fields) Debug(msg string) {
	formatLog(DebugLevel, msg, f)
}

func (f Fields) Info(msg string) {
	formatLog(InfoLevel, msg, f)
}

func (f Fields) Warn(msg string) {
	formatLog(WarnLevel, msg, f)
}

func (f Fields) Error(msg string) {
	formatLog(ErrorLevel, msg, f)
}

func (f Fields) Fatal(msg string) {
	formatLog(FatalLevel, msg, f)
	os.Exit(1)
}

func (f Fields) Panic(msg string) {
	formatLog(PanicLevel, msg, f)
	panic("panic")
}

func (f Fields) Trace(msg string) {
	formatLog(TraceLevel, msg, f)
}
