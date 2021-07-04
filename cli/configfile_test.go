package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_parseConfig(t *testing.T) {
	type args struct {
		s        string
		fileName string
	}
	tests := []struct {
		name     string
		args     args
		want     map[string][]configFileEntry
		wantErr  bool
		errorStr string
	}{
		{
			name: "empty config",
		},
		{
			name: "invalid",
			args: args{
				s: "1234",
			},
			wantErr:  true,
			errorStr: "invalid",
		},
		{
			name: "no config",
			args: args{
				s: "# hello world",
			},
		},
		{
			name: "config w/ value",
			args: args{
				s:        "hello=world",
				fileName: "hi2u",
			},
			want: map[string][]configFileEntry{
				"hello": {
					{
						fileName: "hi2u",
						line:     1,
						key:      "hello",
						hasValue: true,
						value:    "world",
					},
				},
			},
		},
		{
			name: "config w/o value",
			args: args{
				s:        "hello",
				fileName: "hi2u2",
			},
			want: map[string][]configFileEntry{
				"hello": {
					{
						fileName: "hi2u2",
						line:     1,
						key:      "hello",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConfig(strings.NewReader(tt.args.s), tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setConfig(t *testing.T) {
	type args struct {
		entries map[string][]configFileEntry
		cl      ConfigList
	}
	testList := []struct {
		name      string
		args      args
		hasError  bool
		errorStr  string
		checkFunc func(entries map[string][]configFileEntry, cl ConfigList) error
	}{
		{name: "empty"},
		{
			name: "no arg for cfg",
			args: args{
				entries: map[string][]configFileEntry{"hello": {{}}},
			},
			hasError: true,
			errorStr: "unknown config",
		},
		{
			name: "set arg",
			args: args{
				entries: map[string][]configFileEntry{"hello": {{key: "hello", hasValue: true, value: "world"}}},
				cl:      ConfigList{&Config{Name: "hello"}},
			},
			checkFunc: func(config map[string][]configFileEntry, cl ConfigList) error {
				cfg := cl.Get("hello")
				val := cfg.String()
				if val != "world" {
					return fmt.Errorf("world != %s", val)
				}
				return nil
			},
		},
		{
			name: "set config fail",
			args: args{
				entries: map[string][]configFileEntry{"hello": {{key: "hello", hasValue: true, value: "world"}}},
				cl: ConfigList{&Config{Name: "hello", OnSet: func(*Config, []string, ConfigSource) error {
					return errors.New("blah")
				}}},
			},
			hasError: true,
			errorStr: "error setting config",
		},
	}

	for _, currentTest := range testList {
		t.Run(currentTest.name, func(t *testing.T) {
			err := setConfig(currentTest.args.entries, currentTest.args.cl)
			if currentTest.hasError == (err == nil) {
				t.Fatalf("expected error %v, got: %v", currentTest.hasError, err)
			}
			if err != nil && !strings.Contains(err.Error(), currentTest.errorStr) {
				t.Fatalf("expected %q in error: %v", currentTest.errorStr, err)
			}
		})
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

func newTestFile(name string, b []byte) *testFile {
	return &testFile{
		name: name,
		r:    bytes.NewReader(b),
		size: int64(len(b)),
	}
}

type testFile struct {
	name string
	r    io.Reader
	size int64
}

func (tf testFile) Name() string {
	return tf.name
}

func (tf testFile) Size() int64 {
	return tf.size
}

func (tf testFile) Mode() fs.FileMode {
	return 0
}

func (tf testFile) ModTime() time.Time {
	return time.Now()
}

func (tf testFile) IsDir() bool {
	return false
}

func (tf testFile) Sys() interface{} {
	return nil
}

func (tf testFile) Stat() (fs.FileInfo, error) {
	return tf, nil
}

func (tf *testFile) Read(b []byte) (int, error) {
	return tf.r.Read(b)
}

func (tf *testFile) Close() error {
	return nil
}

func Test_loadConfig(t *testing.T) {
	type args struct {
		configFiles    []string
		skipOpenErrors bool
		cl             ConfigList
		sys            fs.FS
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		errorStr  string
		checkFunc func(ConfigList) error
	}{
		{
			name:     "nil",
			args:     args{},
			wantErr:  true,
			errorStr: "no config files",
		},
		{
			name: "empty arg",
			args: args{
				configFiles: []string{},
			},
			wantErr:  true,
			errorStr: "no config files",
		},
		{
			name: "simple default",
			args: args{
				configFiles:    []string{"hello"},
				skipOpenErrors: true,
				sys:            mockFS{},
			},
			wantErr: false,
		},
		{
			name: "error present",
			args: args{
				configFiles: []string{"hello"},
				sys:         mockFS{},
			},
			wantErr:  true,
			errorStr: "error opening config file",
		},
		{
			name: "file error",
			args: args{
				configFiles:    []string{"hello"},
				skipOpenErrors: true,
				sys:            mockFS{file: errorFile{}},
			},
			wantErr:  true,
			errorStr: "errorFile",
		},
		{
			name: "file no error",
			args: args{
				configFiles: []string{"hello"},
				sys:         mockFS{file: emptyFile{}},
			},
			wantErr: false,
		},
		{
			name: "multiple files no error",
			args: args{
				configFiles: []string{"hello1", "hello2"},
				cl:          ConfigList{{Name: "hello"}},
				sys:         mockFS{file: newTestFile("1", []byte("hello=world\n"))},
			},
			wantErr: false,
			checkFunc: func(cl ConfigList) error {
				cfg := cl.Get("hello")
				if cfg == nil {
					return errors.New("no config")
				}
				if !reflect.DeepEqual(cfg.StringSlice(), []string{"world", "world"}) {
					return fmt.Errorf("unexpected value: %+v", cfg.StringSlice())
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loadConfig(tt.args.configFiles, tt.args.skipOpenErrors, tt.args.cl, tt.args.sys)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), tt.errorStr) {
				t.Fatalf("expected %q in error: %v", tt.errorStr, err)
			}
			if err == nil && tt.checkFunc != nil {
				err = tt.checkFunc(tt.args.cl)
				if err != nil {

				}
			}
		})
	}
}
