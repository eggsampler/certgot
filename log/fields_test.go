package log

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type stringerStruct struct {
	S int
}

func (ss stringerStruct) String() string {
	return fmt.Sprintf("%d", ss.S)
}

func TestFields_WithField(t *testing.T) {
	type args struct {
		key   interface{}
		value interface{}
	}
	tests := []struct {
		name  string
		f     Fields
		args  args
		want  Fields
		level Level
	}{
		{
			name: "low",
			f:    Fields{},
			args: args{
				key:   "key",
				value: struct{ value string }{value: "value"},
			},
			want: Fields{stringLogField{
				key:   "key",
				value: "struct { value string }{value:\"value\"}",
			}},
			level: Level(0),
		},
		{
			name: "existing",
			f: Fields{
				stringLogField{
					key:   "foo",
					value: "bar",
				},
			},
			args: args{
				key:   "key",
				value: struct{ value string }{value: "value"},
			},
			want: Fields{
				stringLogField{
					key:   "foo",
					value: "bar",
				},
				stringLogField{
					key:   "key",
					value: "struct { value string }{value:\"value\"}",
				},
			},
			level: Level(0),
		},
		{
			name: "med",
			f:    Fields{},
			args: args{
				key:   "key",
				value: struct{ value string }{value: "value"},
			},
			want: Fields{stringLogField{
				key:   "key",
				value: "{value:value}",
			}},
			level: Level(1),
		},
		{
			name: "high",
			f:    Fields{},
			args: args{
				key:   "key",
				value: struct{ value string }{value: "value"},
			},
			want: Fields{stringLogField{
				key:   "key",
				value: "{value}",
			}},
			level: Level(2),
		},
		{
			name: "stringer",
			f:    Fields{},
			args: args{
				key:   stringerStruct{S: 42},
				value: struct{ value string }{value: "value"},
			},
			want: Fields{stringLogField{
				key:   "42",
				value: "{value}",
			}},
			level: Level(2),
		},
		{
			name: "stringer",
			f:    Fields{},
			args: args{
				key:   struct{ S int }{S: 42},
				value: struct{ value string }{value: "value"},
			},
			want: Fields{stringLogField{
				key:   "{42}",
				value: "{value}",
			}},
			level: Level(2),
		},
		{
			name: "nil",
			f:    Fields{},
			args: args{
				key:   nil,
				value: nil,
			},
			want: Fields{stringLogField{
				key:   "<nil>",
				value: "<nil>",
			}},
			level: Level(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: set the level independent of the package level
			currentLevel = tt.level
			got := tt.f.WithField(tt.args.key, tt.args.value)

			if len(got) != len(tt.want) {
				t.Errorf("expected %d field, got: %d", tt.want, len(got))
			}

			for i := 0; i < len(got); i++ {
				if got[i].Key() != tt.want[i].Key() {
					t.Errorf("key %d mismatch, want: %s, got: %s", i, tt.want[i].Key(), got[i].Key())
				}
				if got[i].Value() != tt.want[i].Value() {
					t.Errorf("value %d mismatch, want: %v, got: %v", i, tt.want[i].Value(), got[i].Value())
				}
			}
		})
	}
}

func TestFields_WithError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name          string
		f             Fields
		args          args
		wantFields    bool
		wantFieldsVal Fields
		want          func(Fields) error
	}{
		{
			name:          "nil",
			f:             nil,
			args:          args{},
			wantFields:    true,
			wantFieldsVal: nil,
		},
		{
			name:          "more nil",
			f:             Fields{},
			args:          args{},
			wantFields:    true,
			wantFieldsVal: Fields{},
		},
		{
			name: "simple",
			f:    Fields{},
			args: args{
				err: errors.New("omg"),
			},
			wantFields: true,
			wantFieldsVal: Fields{
				errorLogField{
					err: errors.New("omg"),
				},
			},
		},
		{
			name: "existing",
			f: Fields{
				stringLogField{
					key:   "foo",
					value: "bar",
				},
			},
			args: args{
				err: errors.New("omg"),
			},
			wantFields: true,
			wantFieldsVal: Fields{
				stringLogField{
					key:   "foo",
					value: "bar",
				},
				errorLogField{
					err: errors.New("omg"),
				},
			},
		},
		{
			name: "less simple",
			f:    Fields{},
			args: args{
				err: Fields{}.CreateError("foo"),
			},
			want: func(fields Fields) error {
				l := len(fields)
				if l != 3 {
					return fmt.Errorf("unexpected field count: %d", l)
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.f.WithError(tt.args.err)
			if tt.wantFields && !reflect.DeepEqual(got, tt.wantFieldsVal) {
				t.Errorf("WithError() = %v, want %v", got, tt.wantFieldsVal)
			}
			if tt.want != nil {
				if err := tt.want(got); err != nil {
					t.Errorf("WithError() = %v, %v", got, err)
				}
			}
		})
	}
}

func TestFields_WithFields(t *testing.T) {
	type args struct {
		fields []interface{}
	}
	tests := []struct {
		name          string
		f             Fields
		args          args
		wantFields    bool
		wantFieldsVal Fields
		want          func(Fields) error
	}{
		{
			name:          "nil",
			f:             nil,
			args:          args{},
			wantFields:    false,
			wantFieldsVal: nil,
			want:          nil,
		},
		{
			name: "simple",
			f:    Fields{},
			args: args{
				fields: []interface{}{
					"hello", "world",
				},
			},
			wantFields: true,
			wantFieldsVal: Fields{
				stringLogField{
					key:   "hello",
					value: "world",
				},
			},
			want: nil,
		},
		{
			name: "existing",
			f: Fields{
				stringLogField{
					key:   "foo",
					value: "bar",
				},
			},
			args: args{
				fields: []interface{}{
					"hello", "world",
				},
			},
			wantFields: true,
			wantFieldsVal: Fields{
				stringLogField{
					key:   "foo",
					value: "bar",
				},
				stringLogField{
					key:   "hello",
					value: "world",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.f.WithFields(tt.args.fields...)
			if tt.wantFields && !reflect.DeepEqual(got, tt.wantFieldsVal) {
				t.Errorf("WithFields() = %v, want %v", got, tt.wantFieldsVal)
			}
			if tt.want != nil {
				if err := tt.want(got); err != nil {
					t.Errorf("WithFields() = %v, %v", got, err)
				}
			}
		})
	}
}
