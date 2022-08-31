package logr

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/mrodden/logr/logger"
)

func Fatal(v ...any)                   { log(logger.ERROR, v...); os.Exit(1) }
func Fatalf(format string, v ...any)   { logf(logger.ERROR, format, v...); os.Exit(1) }
func Error(v ...any)                   { log(logger.ERROR, v...) }
func Errorf(format string, v ...any)   { logf(logger.ERROR, format, v...) }
func Warn(v ...any)                    { log(logger.WARN, v...) }
func Warnf(format string, v ...any)    { logf(logger.WARN, format, v...) }
func Warning(v ...any)                 { log(logger.WARN, v...) }
func Warningf(format string, v ...any) { logf(logger.WARN, format, v...) }
func Info(v ...any)                    { log(logger.INFO, v...) }
func Infof(format string, v ...any)    { logf(logger.INFO, format, v...) }
func Debug(v ...any)                   { log(logger.DEBUG, v...) }
func Debugf(format string, v ...any)   { logf(logger.DEBUG, format, v...) }
func Trace(v ...any)                   { log(logger.TRACE, v...) }
func Tracef(format string, v ...any)   { logf(logger.TRACE, format, v...) }

// CallerPackageName returns the name of the go package of the function at specified depth up the callstack
func CallerPackageName(depth int) string {
	pc, _, _, _ := runtime.Caller(depth)
	fullName := runtime.FuncForPC(pc).Name()
	return cleanupPackageName(fullName)
}

func cleanupPackageName(name string) string {
	idx := strings.LastIndex(name, "/")
	var path string
	if idx >= 0 {
		// remove struct and function names
		lastPart := strings.SplitN(name[idx:], ".", 2)[0]
		path = name[:idx] + lastPart
	} else {
		// no / in string
		// remove struct and function names
		path = strings.SplitN(name, ".", 2)[0]
	}
	return path
}

func pcToName(pc uintptr) string {
	return runtime.FuncForPC(pc).Name()
}

// Caller returns the filename and line number of the function at the specified depth up the callstack
func Caller(depth int) (string, string, uint) {
	pc, file, lineno, _ := runtime.Caller(depth)
	return pcToName(pc), RSplitN(file, "/", 2)[1], uint(lineno)
}

// RSplitN like SplitN but starting from the right
func RSplitN(s, sep string, n int) []string {
	// maintain consistency with SplitN behavior
	if n == 0 {
		return nil
	}
	if n == 1 {
		return []string{s}
	}

	parts := strings.Split(s, sep)
	// maintain consistency with SplitN behavior
	if n < 0 {
		return parts
	}

	// n > 1
	toReturn := []string{strings.Join(parts[0:len(parts)+1-n], sep)}
	toReturn = append(toReturn, parts[len(parts)+1-n:]...)
	return toReturn
}

func log(level logger.Level, v ...any) {
	output(3, level, fmt.Sprint(v...))
}

func logf(level logger.Level, format string, v ...any) {
	output(3, level, fmt.Sprintf(format, v...))
}

func output(depth int, level logger.Level, s string) error {
	if logger.Global() == nil {
		return nil
	}

	pkg, file, line := Caller(depth + 1)
	pkg = cleanupPackageName(pkg)

	rec := logger.NewRecordBuilder().
		Level(level).
		Target(pkg).
		Args(s).
		PackagePath(pkg).
		File(file).
		Line(line).
		Build()

	logger.Global().Log(rec)
	return nil
}
