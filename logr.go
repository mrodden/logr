package logr

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-isatty"
)

type Logger struct {
	mu         sync.Mutex
	out        io.Writer
	forceColor bool
}

func New(out io.Writer) *Logger {
	return &Logger{out: out}
}

var std = New(os.Stderr)

func Default() *Logger { return std }

func Fatal(v ...any)                   { std.log(CRITICAL, v...); os.Exit(1) }
func Fatalf(format string, v ...any)   { std.logf(CRITICAL, format, v...); os.Exit(1) }
func Error(v ...any)                   { std.log(ERROR, v...) }
func Errorf(format string, v ...any)   { std.logf(ERROR, format, v...) }
func Warn(v ...any)                    { std.log(WARN, v...) }
func Warnf(format string, v ...any)    { std.logf(WARN, format, v...) }
func Warning(v ...any)                 { std.log(WARN, v...) }
func Warningf(format string, v ...any) { std.logf(WARN, format, v...) }
func Info(v ...any)                    { std.log(INFO, v...) }
func Infof(format string, v ...any)    { std.logf(INFO, format, v...) }
func Debug(v ...any)                   { std.log(DEBUG, v...) }
func Debugf(format string, v ...any)   { std.logf(DEBUG, format, v...) }

type LogLevel uint8

const (
	CRITICAL LogLevel = iota
	ERROR
	WARN
	INFO
	DEBUG
)

const (
	ISO = "2006-01-02T15:04:05.999Z"
)

func (l *Logger) Fatal(v ...any)                   { l.log(CRITICAL, v...); os.Exit(1) }
func (l *Logger) Fatalf(format string, v ...any)   { l.logf(CRITICAL, format, v...); os.Exit(1) }
func (l *Logger) Error(v ...any)                   { l.log(ERROR, v...) }
func (l *Logger) Errorf(format string, v ...any)   { l.logf(ERROR, format, v...) }
func (l *Logger) Warn(v ...any)                    { l.log(WARN, v...) }
func (l *Logger) Warnf(format string, v ...any)    { l.logf(WARN, format, v...) }
func (l *Logger) Warning(v ...any)                 { l.log(WARN, v...) }
func (l *Logger) Warningf(format string, v ...any) { l.logf(WARN, format, v...) }
func (l *Logger) Info(v ...any)                    { l.log(INFO, v...) }
func (l *Logger) Infof(format string, v ...any)    { l.logf(INFO, format, v...) }
func (l *Logger) Debug(v ...any)                   { l.log(DEBUG, v...) }
func (l *Logger) Debugf(format string, v ...any)   { l.logf(DEBUG, format, v...) }

func (l *Logger) ForceColor(force bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.forceColor = force
}

func ltoa(level LogLevel) string {
	switch level {
	case CRITICAL:
		return "CRITICAL"
	case ERROR:
		return "ERROR"
	case WARN:
		return "WARN"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	default:
		return "UNDEF"
	}
}

// CallerPackageName returns the name of the go package of the function at specified depth up the callstack
func CallerPackageName(depth int) string {
	pc, _, _, _ := runtime.Caller(depth)
	fullName := runtime.FuncForPC(pc).Name()
	pkgPath := RSplitN(fullName, ".", 2)[0]
	pkg := RSplitN(pkgPath, "/", 2)[1]
	return pkg
}

// Caller returns the filename and line number of the function at the specified depth up the callstack
func Caller(depth int) string {
	_, file, lineno, _ := runtime.Caller(depth)
	return RSplitN(file, "/", 2)[1] + ":" + strconv.Itoa(lineno)
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

func GoRoutineID() string {
	return string(bytes.Fields(debug.Stack())[1])
}

func (l *Logger) log(level LogLevel, v ...any) {
	l.output(level, fmt.Sprint(v...))
}

func (l *Logger) logf(level LogLevel, format string, v ...any) {
	l.output(level, fmt.Sprintf(format, v...))
}

func shouldColorize(out io.Writer) bool {
	f, ok := out.(*os.File)
	if runtime.GOOS == "windows" {
		ok = false
	}
	return ok && isatty.IsTerminal(f.Fd())
}

func (l *Logger) output(level LogLevel, s string) error {
	// time thread component condition
	//lf := "%-24s %02v %s %5s: %s\n"

	dimmed := Style{}
	if shouldColorize(l.out) || l.forceColor {
		dimmed = Style{}.Dimmed()
	}

	ts := dimmed.Paint(fmt.Sprintf("%-24s ", time.Now().UTC().Format(ISO)))

	thr := dimmed.Paint(fmt.Sprintf("%02v ", GoRoutineID()))

	cond := fmt.Sprintf("%5s ", ltoa(level))
	if shouldColorize(l.out) || l.forceColor {
		cond = colorize(level, cond)
	}

	comp := dimmed.Paint(Caller(4))

	buf := make([]byte, 0)
	buf = append(buf, ts...)
	buf = append(buf, thr...)
	buf = append(buf, cond...)
	buf = append(buf, comp...)
	buf = append(buf, dimmed.Paint(":")...)
	buf = append(buf, ' ')
	buf = append(buf, s...)

	if len(s) == 0 || s[len(s)-1] != '\n' {
		buf = append(buf, '\n')
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	_, err := l.out.Write(buf)
	return err
}

func colorize(level LogLevel, s string) string {
	switch level {
	case DEBUG:
		return Blue.Paint(s)
	case INFO:
		return Green.Paint(s)
	case WARN:
		return Yellow.Paint(s)
	case ERROR:
		return Red.Paint(s)
	case CRITICAL:
		return Red.Paint(s)
	default:
		return Green.Paint(s)
	}
}
