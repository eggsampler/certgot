package log

import (
	"errors"
	"fmt"
)

// CreateError creates an error that can be displayed for just this error message
// or alternatively typed to the error value with the log fields extracted
func (f Fields) CreateError(msg string) error {
	return errorWithLogFields{
		errorMessage: msg,
		logFields:    addErrorFields(f),
	}
}

// CreateErrorF is a helper function for CreateError
// or alternatively typed to the error value with the log fields extracted
func (f Fields) CreateErrorF(format string, args ...interface{}) error {
	return f.CreateError(fmt.Sprintf(format, args...))
}

// GetError returns the message passed to CreateError as well as the error used in log.WithError
func GetError(err error) (string, string) {
	if err == nil {
		return "", ""
	}
	var e errorWithLogFields
	if errors.As(err, &e) {
		for _, f := range e.logFields {
			if f.Key() == "error" {
				return e.errorMessage, f.Value()
			}
		}
		return e.errorMessage, ""
	}
	return err.Error(), ""
}

func CreateError(msg string) error {
	return Fields{}.CreateError(msg)
}

func addErrorFields(f Fields) Fields {
	addedFields := []logField{
		stringLogField{"time", getTime()},
		stringLogField{"source", getSource()},
	}
	return append(f, addedFields...)
}

type errorWithLogFields struct {
	errorMessage string
	logFields    Fields
}

func (e errorWithLogFields) Error() string {
	return e.errorMessage
}

func getErrorFields(err error) Fields {
	var e errorWithLogFields
	if errors.As(err, &e) {
		return e.logFields
	}
	return nil
}
