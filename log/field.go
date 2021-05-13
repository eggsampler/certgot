package log

type logField interface {
	Key() string
	Value() string
}

type errorLogField struct {
	err error
}

func (e errorLogField) Key() string {
	return "error"
}

func (e errorLogField) Value() string {
	return e.err.Error()
}

type stringLogField struct {
	key, value string
}

func (s stringLogField) Key() string {
	return s.key
}

func (s stringLogField) Value() string {
	return s.value
}
