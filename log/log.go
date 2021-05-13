package log

// This file is intended to be used to begin the method chaining for the logging
// That is, essentially a drop in replacement for the default go "log" package

func WithError(err error) Fields {
	return Fields{}.WithError(err)
}

func WithField(key string, value interface{}) Fields {
	return Fields{}.WithField(key, value)
}

func WithFields(fields ...interface{}) Fields {
	return Fields{}.WithFields(fields...)
}

func Printf(format string, a ...interface{}) {
	Fields{}.Printf(format, a...)
}

func Log(level Level, message string) {
	Fields{}.Log(level, message)
}

func Debug(message string) {
	Fields{}.Debug(message)
}

func Info(message string) {
	Fields{}.Info(message)
}

func Warn(message string) {
	Fields{}.Warn(message)
}

func Error(message string) {
	Fields{}.Error(message)
}

func Fatal(message string) {
	Fields{}.Fatal(message)
}

func Panic(message string) {
	Fields{}.Panic(message)
}

func Trace(message string) {
	Fields{}.Trace(message)
}
