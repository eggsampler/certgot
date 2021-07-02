package cli

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_parseArguments(t *testing.T) {
	type args struct {
		argsToParse []string
		ctx         *Context
		fl          FlagList
		cl          CommandList
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		errStr    string
		checkFunc func(ctx *Context) error
	}{
		{name: "empty"},
		{
			name: "binary only",
			args: args{argsToParse: []string{"bin"}},
		},
		{
			name: "invalid flag",
			args: args{
				argsToParse: []string{"bin", "-"},
			},
			wantErr: true,
			errStr:  "invalid flag",
		},
		{
			name: "unknown short flag",
			args: args{
				argsToParse: []string{"bin", "-d"},
			},
			wantErr: true,
			errStr:  "unknown flag",
		},
		{
			name: "unknown long flag",
			args: args{
				argsToParse: []string{"bin", "--d"},
			},
			wantErr: true,
			errStr:  "unknown flag",
		},
		{
			name: "bad repeated short flag",
			args: args{
				argsToParse: []string{"bin", "-dddd"},
				ctx:         &Context{},
				fl:          FlagList{&Flag{Name: "d"}},
			},
			wantErr: true,
			errStr:  "repeated short flags",
		},
		{
			name: "bad multiple",
			args: args{
				argsToParse: []string{"bin", "-d", "-d"},
				ctx:         &Context{},
				fl:          FlagList{&Flag{Name: "d"}},
			},
			wantErr: true,
			errStr:  "multiple",
		},
		{
			name: "simple ok",
			args: args{
				argsToParse: []string{"bin", "-d"},
				ctx:         &Context{},
				fl:          FlagList{&Flag{Name: "d"}},
			},
			checkFunc: func(ctx *Context) error {
				if len(ctx.Flags) != 1 {
					return fmt.Errorf("bad flag count: %d", len(ctx.Flags))
				}
				return nil
			},
		},
		{
			name: "requires value",
			args: args{
				argsToParse: []string{"bin", "-d"},
				ctx:         &Context{},
				fl:          FlagList{&Flag{Name: "d", TakesValue: true}},
			},
			wantErr: true,
			errStr:  "requires value",
		},
		{
			name: "inline doesn't take value",
			args: args{
				argsToParse: []string{"bin", "-d=1"},
				ctx:         &Context{},
				fl:          FlagList{&Flag{Name: "d"}},
			},
			wantErr: true,
			errStr:  "doesn't take a value",
		},
		{
			name: "inline ok value",
			args: args{
				argsToParse: []string{"bin", "-d=1"},
				ctx:         &Context{},
				fl:          FlagList{&Flag{Name: "d", TakesValue: true}},
			},
			checkFunc: func(ctx *Context) error {
				f := ctx.Flags.Get("d")
				if f == nil {
					return errors.New("ctx didn't include flag")
				}
				if len(f.values) == 0 {
					return errors.New("flag didn't include any values")
				}
				if f.values[0] != "1" {
					return fmt.Errorf("unexpected flag value: %+v", f.values)
				}
				return nil
			},
		},
		{
			name: "ok value",
			args: args{
				argsToParse: []string{"bin", "-d", "1"},
				ctx:         &Context{},
				fl:          FlagList{&Flag{Name: "D", TakesValue: true}},
			},
			checkFunc: func(ctx *Context) error {
				f := ctx.Flags.Get("d")
				if f == nil {
					return errors.New("ctx didn't include flag")
				}
				if len(f.values) == 0 {
					return errors.New("flag didn't include any values")
				}
				if f.values[0] != "1" {
					return fmt.Errorf("unexpected flag value: %+v", f.values)
				}
				return nil
			},
		},
		{
			name: "invalid command",
			args: args{
				argsToParse: []string{"bin", "cmd"},
				ctx:         &Context{},
			},
			wantErr: true,
			errStr:  "invalid command",
		},
		{
			name: "ok args",
			args: args{
				argsToParse: []string{"bin", "cmd", "args"},
				ctx:         &Context{},
				cl:          CommandList{&Command{Name: "CMD"}},
			},
			checkFunc: func(ctx *Context) error {
				if ctx.Command == nil {
					return errors.New("no command")
				}
				if ctx.Command.Name != "CMD" {
					return fmt.Errorf("invalid cmd: %s", ctx.Command.Name)
				}
				if !reflect.DeepEqual(ctx.ExtraArguments, []string{"args"}) {
					return fmt.Errorf("invalid args: %+v", ctx.ExtraArguments)
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseArguments(tt.args.argsToParse, tt.args.ctx, tt.args.fl, tt.args.cl)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArguments() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), tt.errStr) {
				t.Errorf("expected %q in error: %v", tt.errStr, err)
			}
			if err == nil && tt.checkFunc != nil {
				err = tt.checkFunc(tt.args.ctx)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func Test_extractFlag(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want []string
	}{
		{
			name: "empty",
		},
		{
			name: "invalid short",
			arg:  "-",
		},
		{
			name: "invalid long",
			arg:  "--",
		},
		{
			name: "ok short",
			arg:  "-a",
			want: []string{"-a", "a", ""},
		},
		{
			name: "bad short",
			arg:  "-ab",
		},
		{
			name: "bad short no value",
			arg:  "-a=",
		},
		{
			name: "bad short value",
			arg:  "-ab=a",
		},
		{
			name: "bad long",
			arg:  "---",
		},
		{
			name: "long no value",
			arg:  "--abc=",
		},
		{
			name: "ok long value",
			arg:  "--abc=asd",
			want: []string{"--abc=asd", "abc", "asd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractFlag(tt.arg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}
