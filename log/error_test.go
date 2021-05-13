package log

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestFields_CreateError(t *testing.T) {
	err := Fields{}.CreateError("hello, world")
	if err.Error() != "hello, world" {
		t.Error("wrong error text")
	}
	errLog, ok := err.(errorWithLogFields)
	if !ok {
		t.Error("no fields")
	}
	if len(errLog.logFields) != 2 {
		t.Error("no source/time")
	}
	if errLog.errorMessage != "hello, world" {
		t.Error("wrong error field text")
	}
}

func Test_errorWithLogFields_Error(t *testing.T) {
	type fields struct {
		errorMessage string
		logFields    Fields
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "exist",
			fields: fields{
				errorMessage: "hello",
			},
			want: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := errorWithLogFields{
				errorMessage: tt.fields.errorMessage,
				logFields:    tt.fields.logFields,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getErrorFields(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want Fields
	}{
		{},
		{
			name: "nil",
			args: args{nil},
		},
		{
			name: "simple error",
			args: args{errors.New("simple")},
		},
		{
			name: "no fields",
			args: args{Fields{}.CreateError("no fields error")},
			want: addErrorFields(Fields{}),
		},
		{
			name: "some fields",
			args: args{WithField("hello", "world").CreateError("some fields error")},
			want: addErrorFields(WithField("hello", "world")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getErrorFields(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getErrorFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addErrorFields(t *testing.T) {
	type args struct {
		f Fields
	}
	tests := []struct {
		name  string
		args  args
		check func(got Fields) error
	}{
		{
			name: "simple",
			args: args{Fields{stringLogField{"hello", "world"}}},
			check: func(got Fields) error {
				if len(got) != 3 {
					return fmt.Errorf("unexpected length, want: 3, got: %d", len(got))
				}
				allowedKeys := map[string]bool{"hello": true, "time": true, "source": true}
				for _, v := range got {
					if !allowedKeys[v.Key()] {
						return fmt.Errorf("unknown key: %s", v.Key())
					}
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addErrorFields(tt.args.f)
			if tt.check != nil {
				if err := tt.check(got); err != nil {
					t.Errorf("addErrorFields() = %v: %v", got, err)
				}
			}
		})
	}
}

func TestGetError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "nil",
			args: args{
				err: nil,
			},
			want:  "",
			want1: "",
		},
		{
			name: "simple",
			args: args{
				err: CreateError("hello, world"),
			},
			want:  "hello, world",
			want1: "",
		},
		{
			name: "standard error",
			args: args{
				err: errors.New("ok boomer"),
			},
			want:  "ok boomer",
			want1: "",
		},
		{
			name: "error",
			args: args{
				err: WithError(errors.New("foo")).CreateError("bar"),
			},
			want:  "bar",
			want1: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetError(tt.args.err)
			if got != tt.want {
				t.Errorf("GetError() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetError() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
