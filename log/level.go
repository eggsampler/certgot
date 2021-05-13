package log

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel

	TraceLevel Level = -1

	MinLevel = TraceLevel
	MaxLevel = PanicLevel
)

var (
	currentLevel = Level(0)
)

func SetLevel(level Level) {
	currentLevel = level
}

func GetLevel() Level {
	return currentLevel
}
