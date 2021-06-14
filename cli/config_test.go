package cli

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func Test_parseConfig(t *testing.T) {
	emptyCfg := func(cfg []configEntry) error {
		if len(cfg) != 0 {
			return fmt.Errorf("unexpected config: %v", cfg)
		}
		return nil
	}

	testList := []struct {
		testName     string
		configString string
		fileName     string
		hasError     bool
		errorStr     string
		checkFunc    func(cfg []configEntry) error
	}{
		{
			testName:  "empty config",
			checkFunc: emptyCfg,
		},
		{
			testName:     "invalid",
			configString: "1234",
			hasError:     true,
			errorStr:     "invalid",
		},
		{
			testName:     "no config",
			configString: "# hello world",
			checkFunc:    emptyCfg,
		},
		{
			testName:     "config w/ value",
			configString: "hello=world",
			fileName:     "hi2u",
			checkFunc: func(cfg []configEntry) error {
				if len(cfg) != 1 {
					return errors.New("not 1 config")
				}
				entry := cfg[0]
				otherEntry := configEntry{
					fileName: "hi2u",
					line:     1,
					key:      "hello",
					hasValue: true,
					value:    "world",
				}
				if !reflect.DeepEqual(entry, otherEntry) {
					return fmt.Errorf("entry mismatch %+v != %+v", entry, otherEntry)
				}
				return nil
			},
		},
		{
			testName:     "config w/o value",
			configString: "hello",
			fileName:     "hi2u2",
			checkFunc: func(cfg []configEntry) error {
				if len(cfg) != 1 {
					return errors.New("not 1 config")
				}
				entry := cfg[0]
				otherEntry := configEntry{
					fileName: "hi2u2",
					line:     1,
					key:      "hello",
				}
				if !reflect.DeepEqual(entry, otherEntry) {
					return fmt.Errorf("entry mismatch %+v != %+v", entry, otherEntry)
				}
				return nil
			},
		},
	}

	for _, currentTest := range testList {
		cfg, err := parseConfig(strings.NewReader(currentTest.configString), currentTest.fileName)
		if currentTest.hasError == (err == nil) {
			t.Fatalf("%q: expected error %v, got: %v", currentTest.testName, currentTest.hasError, err)
		}
		if err != nil && !strings.Contains(err.Error(), currentTest.errorStr) {
			t.Fatalf("test %q: expected %q in error: %v", currentTest.testName, currentTest.errorStr, err)
		}
		if currentTest.checkFunc != nil {
			if err := currentTest.checkFunc(cfg); err != nil {
				t.Fatalf("%q: unexpected error %v", currentTest.testName, err)
			}
		}
	}
}

func Test_setConfig(t *testing.T) {
	testList := []struct {
		testName  string
		config    []configEntry
		args      map[string]*Argument
		hasError  bool
		errorStr  string
		checkFunc func(config map[string]configEntry, args map[string]*Argument) error
	}{
		{testName: "empty"},
		{
			testName: "no arg for cfg",
			config:   []configEntry{{}},
			hasError: true,
			errorStr: "unknown argument",
		},
		{
			testName: "set arg",
			config:   []configEntry{{key: "hello", hasValue: true, value: "world"}},
			args: map[string]*Argument{"hello": {
				TakesValue: true,
			}},
			checkFunc: func(config map[string]configEntry, args map[string]*Argument) error {
				arg := args["hello"]
				val := arg.String()
				if val != "world" {
					return fmt.Errorf("world != %s", val)
				}
				return nil
			},
		},
		{
			testName: "set arg fail",
			config:   []configEntry{{key: "hello", hasValue: true, value: "world"}},
			args:     map[string]*Argument{"hello": {}},
			hasError: true,
			errorStr: "error setting arg",
		},
	}

	for _, currentTest := range testList {
		err := setConfig(currentTest.config, currentTest.args)
		if currentTest.hasError == (err == nil) {
			t.Fatalf("%q: expected error %v, got: %v", currentTest.testName, currentTest.hasError, err)
		}
		if err != nil && !strings.Contains(err.Error(), currentTest.errorStr) {
			t.Fatalf("test %q: expected %q in error: %v", currentTest.testName, currentTest.errorStr, err)
		}
	}
}

var (
	funcNoEnv  = func(string) string { return "" }
	funcNoUser = func() (*user.User, error) { return nil, errors.New("no user") }
)

func Test_parsePath(t *testing.T) {
	type args struct {
		path     string
		envFunc  func(string) string
		userFunc func() (*user.User, error)
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				path:     "",
				envFunc:  funcNoEnv,
				userFunc: funcNoUser,
			},
			want: ".",
		},
		{
			name: "simple",
			args: args{
				path:     "~",
				envFunc:  funcNoEnv,
				userFunc: funcNoUser,
			},
			want: "~",
		},
		{
			name: "xdg 1",
			args: args{
				path: "~",
				envFunc: func(s string) string {
					if s == "XDG_CONFIG_HOME" {
						return "XDG_CONFIG_HOME"
					}
					return ""
				},
				userFunc: funcNoUser,
			},
			want: "XDG_CONFIG_HOME",
		},
		{
			name: "xdg 2",
			args: args{
				path: filepath.Join("~", "hello"),
				envFunc: func(s string) string {
					if s == "XDG_CONFIG_HOME" {
						return "XDG_CONFIG_HOME"
					}
					return ""
				},
				userFunc: funcNoUser,
			},
			want: filepath.Join("XDG_CONFIG_HOME", "hello"),
		},
		{
			name: "user home",
			args: args{
				path:    "~",
				envFunc: funcNoEnv,
				userFunc: func() (*user.User, error) {
					return &user.User{
						HomeDir: "USER_HOME_DIR",
					}, nil
				},
			},
			want: "USER_HOME_DIR",
		},
		{
			name: "home",
			args: args{
				path: "~",
				envFunc: func(s string) string {
					if s == "HOME" {
						return "HOME"
					}
					return ""
				},
				userFunc: funcNoUser,
			},
			want: "HOME",
		},
		{
			name: "home drive/path",
			args: args{
				path: "~",
				envFunc: func(s string) string {
					if s == "HomeDrive" {
						return "HomeDrive"
					}
					if s == "HomePath" {
						return "HomePath"
					}
					return ""
				},
				userFunc: funcNoUser,
			},
			want: filepath.Join("HomeDrive", "HomePath"),
		},
		{
			name: "userprofile",
			args: args{
				path: "~",
				envFunc: func(s string) string {
					if s == "UserProfile" {
						return "UserProfile"
					}
					return ""
				},
				userFunc: funcNoUser,
			},
			want: "UserProfile",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePath(tt.args.path, tt.args.envFunc, tt.args.userFunc); got != tt.want {
				t.Errorf("parsePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockFS struct {
	file fs.File
}

func (m mockFS) Open(name string) (fs.File, error) {
	if m.file == nil {
		return nil, fs.ErrNotExist
	}
	return m.file, nil
}

type errorFile struct{}

func (errorFile) Stat() (fs.FileInfo, error) {
	return nil, errors.New("errorFile")
}

func (errorFile) Read([]byte) (int, error) {
	return 0, errors.New("errorFile")
}

func (errorFile) Close() error {
	return errors.New("errorFile")
}

type emptyFile struct{}

func (emptyFile) Stat() (fs.FileInfo, error) {
	return nil, nil
}

func (emptyFile) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (emptyFile) Close() error {
	return nil
}

func Test_loadConfig(t *testing.T) {
	type args struct {
		app     *App
		cfgFile *Argument
		sys     fs.FS
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		errorStr string
	}{
		{
			name:     "nil",
			args:     args{},
			wantErr:  true,
			errorStr: "no config",
		},
		{
			name: "empty arg",
			args: args{
				app:     &App{},
				cfgFile: &Argument{},
			},
			wantErr: false,
		},
		{
			name: "simple default",
			args: args{
				app: &App{},
				cfgFile: &Argument{
					DefaultValue: SimpleValue{Value: []string{"hello"}},
				},
				sys: mockFS{},
			},
			wantErr: false,
		},
		{
			name: "error present",
			args: args{
				app: &App{},
				cfgFile: &Argument{
					DefaultValue: SimpleValue{Value: []string{"hello"}},
					isPresent:    true,
				},
				sys: mockFS{},
			},
			wantErr:  true,
			errorStr: "error opening config file",
		},
		{
			name: "file error",
			args: args{
				app: &App{},
				cfgFile: &Argument{
					DefaultValue: SimpleValue{Value: []string{"hello"}},
					isPresent:    true,
				},
				sys: mockFS{file: errorFile{}},
			},
			wantErr:  true,
			errorStr: "errorFile",
		},
		{
			name: "file no error",
			args: args{
				app: &App{},
				cfgFile: &Argument{
					DefaultValue: SimpleValue{Value: []string{"hello"}},
					isPresent:    true,
				},
				sys: mockFS{file: emptyFile{}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loadConfig(tt.args.app, tt.args.cfgFile, tt.args.sys)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), tt.errorStr) {
				t.Fatalf("expected %q in error: %v", tt.errorStr, err)
			}
		})
	}
}
