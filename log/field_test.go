package log

import (
	"errors"
	"testing"
)

func Test_errorLogField(t *testing.T) {
	testErr := errors.New("hello, world")
	elf := errorLogField{testErr}
	if elf.Key() != "error" {
		t.Errorf("unknown key: %s", elf.Key())
	}
	if testErr.Error() != elf.Value() {
		t.Errorf("unknown value: %s", elf.Value())
	}
}

func Test_stringLogField(t *testing.T) {
	slf := stringLogField{
		key:   "hello",
		value: "world",
	}
	if slf.Key() != "hello" {
		t.Errorf("unknown key: %s", slf.Key())
	}
	if slf.Value() != "world" {
		t.Errorf("unknown value: %s", slf.Value())
	}
}
