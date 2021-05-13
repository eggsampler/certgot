package log

import (
	"regexp"
	"runtime/debug"
	"testing"
)

func Test_getTime(t *testing.T) {
	tme := getTime()
	reg := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+\d{2}:\d{2}`)
	if !reg.MatchString(tme) {
		t.Fatalf("bad time: %s", tme)
	}
}

type testDebugProvider struct {
	bi  *debug.BuildInfo
	res bool
}

func (tdp testDebugProvider) ReadBuildInfo() (*debug.BuildInfo, bool) {
	return tdp.bi, tdp.res
}

type callerResult struct {
	pc   uintptr
	file string
	line int
	ok   bool
}

type nameResult struct {
	name string
}

func (nr nameResult) Name() string {
	return nr.name
}

type testRuntimeProvider struct {
	callers map[int]callerResult
	funcs   map[uintptr]nameResult
}

func (trp testRuntimeProvider) Caller(skip int) (pc uintptr, file string, line int, ok bool) {
	if trp.callers == nil {
		return 0, "", 0, false
	}
	return trp.callers[skip].pc, trp.callers[skip].file, trp.callers[skip].line, trp.callers[skip].ok
}

func (trp testRuntimeProvider) FuncForPC(pc uintptr) nameProvider {
	if trp.funcs == nil {
		return nil
	}
	return trp.funcs[pc]
}

func Test_getSourceProvider(t *testing.T) {
	type args struct {
		bip debugProvider
		cp  runtimeProvider
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no buildinfo",
			args: args{
				bip: nil,
				cp:  nil,
			},
			want: "no_buildinfo",
		},
		{
			name: "no caller",
			args: args{
				bip: testDebugProvider{nil, false},
				cp:  nil,
			},
			want: "no_caller",
		},
		{
			name: "unknown build",
			args: args{
				bip: testDebugProvider{nil, false},
				cp:  testRuntimeProvider{},
			},
			want: "unknown_build",
		},
		{
			name: "unknown no build",
			args: args{
				bip: testDebugProvider{nil, true},
				cp:  testRuntimeProvider{},
			},
			want: "unknown_nobuild",
		},
		{
			name: "unknown caller",
			args: args{
				bip: testDebugProvider{&debug.BuildInfo{}, true},
				cp:  testRuntimeProvider{},
			},
			want: "unknown_caller",
		},
		{
			name: "basic call",
			args: args{
				bip: testDebugProvider{&debug.BuildInfo{}, true},
				cp: testRuntimeProvider{
					callers: map[int]callerResult{
						1: {
							0, "hello.go", 0, true,
						},
					}},
			},
			want: "hello.go:0 0",
		},
		{
			name: "slightly basic call",
			args: args{
				bip: testDebugProvider{&debug.BuildInfo{}, true},
				cp: testRuntimeProvider{
					callers: map[int]callerResult{
						1: {
							0, "", 0, true,
						},
						2: {
							0, "hello.go", 0, true,
						},
					}},
			},
			want: "hello.go:0 0",
		},
		{
			name: "less basic call",
			args: args{
				bip: testDebugProvider{&debug.BuildInfo{
					Main: debug.Module{
						Path: "/world/",
					},
				}, true},
				cp: testRuntimeProvider{
					callers: map[int]callerResult{
						1: {
							0, "/world/log/hello.go", 0, true,
						},
						2: {
							0, "/foo/bar.go", 0, true,
						},
					}},
			},
			want: "/foo/bar.go:0 0",
		},
		{
			name: "more less basic call",
			args: args{
				bip: testDebugProvider{&debug.BuildInfo{
					Main: debug.Module{
						Path: "/world/",
					},
				}, true},
				cp: testRuntimeProvider{
					callers: map[int]callerResult{
						1: {
							0, "/world/log/hello.go", 0, true,
						},
						2: {
							0, "/foo/bar.go", 0, true,
						},
					},
					funcs: map[uintptr]nameResult{
						0: {name: "name"},
					}},
			},
			want: "/foo/bar.go:0 name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSourceProvider(tt.args.bip, tt.args.cp); got != tt.want {
				t.Errorf("getSourceProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
