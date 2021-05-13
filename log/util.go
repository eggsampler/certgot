package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

type debugProvider interface {
	ReadBuildInfo() (*debug.BuildInfo, bool)
}

type runtimeProvider interface {
	Caller(skip int) (pc uintptr, file string, line int, ok bool)
	FuncForPC(pc uintptr) nameProvider
}

type nameProvider interface {
	Name() string
}

type normalProvider struct{}

func (normalProvider) ReadBuildInfo() (*debug.BuildInfo, bool) {
	return debug.ReadBuildInfo()
}

func (normalProvider) Caller(skip int) (pc uintptr, file string, line int, ok bool) {
	return runtime.Caller(skip)
}

func (normalProvider) FuncForPC(pc uintptr) nameProvider {
	return runtime.FuncForPC(pc)
}

var levelNames = map[Level]string{
	DebugLevel: "DBG",
	InfoLevel:  "NFO",
	WarnLevel:  "WNR",
	ErrorLevel: "ERR",
	FatalLevel: "FTL",
	PanicLevel: "PNC",
	TraceLevel: "TRC",
}

// getSource finds the first caller outside of the log package
func getSource() string {
	rp := normalProvider{}
	return getSourceProvider(rp, rp)
}

func getSourceProvider(dp debugProvider, rp runtimeProvider) string {
	if dp == nil {
		return "no_buildinfo"
	}
	if rp == nil {
		return "no_caller"
	}
	mod, ok := dp.ReadBuildInfo()
	if !ok {
		return "unknown_build"
	}
	if mod == nil {
		return "unknown_nobuild"
	}
	pkg := mod.Main.Path
	pkgLog := filepath.Join(pkg, "log")
	var pc uintptr
	var file string
	var line int
	for i := 1; i < 10; i++ {
		callerPc, callerFile, callerLine, ok := rp.Caller(i)
		if !ok {
			break
		}
		if callerFile == "" {
			continue
		}
		pc, file, line = callerPc, callerFile, callerLine
		// fmt.Printf("caller skip:%d, pc:%x, file:%q, line:%d\n", i, pc, file, line)
		if strings.Index(callerFile, pkgLog) < 0 {
			break
		}
	}
	if file == "" {
		return "unknown_caller"
	}
	funcName := fmt.Sprintf("%x", pc)
	f := rp.FuncForPC(pc)
	if f != nil {
		funcName = f.Name()
	}
	funcName = filepath.Base(funcName)
	n := strings.Index(file, pkg)
	if n >= 0 {
		return fmt.Sprintf("%s:%d %s", file[n:], line, funcName)
	}
	return fmt.Sprintf("%s:%d %s", file, line, funcName)
}

func getTime() string {
	return time.Now().Format(time.RFC3339) // "2006-01-02 15:04:05 -0700")
}
