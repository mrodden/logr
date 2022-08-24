package env_logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-isatty"

	"github.com/mrodden/logr/logger"
)

const (
	ISO = "2006-01-02T15:04:05.999Z"
)

func TryInit() error {
	l := FromDefaultEnv()

	// set global logr Logger
	logger.SetDefaultLogger(l)

	return nil
}

type envLogger struct {
	directives []Directive

	mu         sync.Mutex
	out        io.Writer
	forceColor bool
}

func FromDefaultEnv() *envLogger {
	return FromEnv(os.Getenv("GO_LOG"))
}

func FromEnv(env string) *envLogger {
	directives := parse(env)

	if len(directives) == 0 {
		directives = append(directives, Directive{Name: "", Level: logger.ERROR})
	}

	return &envLogger{directives: directives, out: os.Stderr}
}

func (l *envLogger) ForceColor(force bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.forceColor = force
}

func shouldColorize(out io.Writer) bool {
	f, ok := out.(*os.File)
	if runtime.GOOS == "windows" {
		ok = false
	}
	return ok && isatty.IsTerminal(f.Fd())
}

func GoRoutineID() string {
	return string(bytes.Fields(debug.Stack())[1])
}

func (l *envLogger) Enabled(metadata logger.Metadata) bool {

	for _, dir := range l.directives {
		if dir.Name != "" && !strings.HasPrefix(metadata.Target(), dir.Name) {
			continue
		}
		return metadata.Level() <= dir.Level
	}

	return false
}

func (l *envLogger) Log(record logger.Record) {
	if l.Enabled(*record.Metadata()) {

		// time thread condition component
		//lf := "%-24s %02v %5s %s: %s\n"

		level := record.Metadata().Level()
		comp := record.Metadata().Target()

		dimmed := Style{}
		if shouldColorize(l.out) || l.forceColor {
			dimmed = Style{}.Dimmed()
		}

		ts := dimmed.Paint(fmt.Sprintf("%-24s ", time.Now().UTC().Format(ISO)))

		thr := dimmed.Paint(fmt.Sprintf("%02v ", GoRoutineID()))

		cond := fmt.Sprintf("%5s ", logger.Ltoa(level))
		if shouldColorize(l.out) || l.forceColor {
			cond = colorize(level, cond)
		}

		s := fmt.Sprint(record.Args()...)

		comp = dimmed.Paint(comp)

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

		_, _ = l.out.Write(buf)
	}
}

func (l *envLogger) Flush() {}

type Directive struct {
	Name  string
	Level logger.Level
}

// ByName fulfills the 'sort.Interface' for sorting Directives by name
type ByName []Directive

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return n[i].Name < n[j].Name }

func parseLevel(s string) (logger.Level, error) {
	l := strings.ToUpper(strings.TrimSpace(s))

	switch l {
	case "TRACE":
		return logger.TRACE, nil
	case "DEBUG":
		return logger.DEBUG, nil
	case "INFO":
		return logger.INFO, nil
	case "WARN":
		return logger.WARN, nil
	case "ERROR":
		return logger.ERROR, nil
	default:
		return logger.TRACE, fmt.Errorf("Unknown log level: %v", l)
	}
}

func parse(ds string) []Directive {
	// comma separated
	// LEVEL
	// TARGET=LEVEL
	// TARGET

	var directives []Directive

	if len(ds) == 0 {
		return directives
	}

	for _, dir := range strings.Split(ds, ",") {

		parts := strings.SplitN(dir, "=", 2)
		t := strings.TrimSpace(parts[0])

		var directive Directive

		if len(parts) == 1 {
			if lv, err := parseLevel(t); err == nil {
				// LEVEL
				directive = Directive{Name: "", Level: lv}
			} else {
				// TARGET
				directive = Directive{Name: t, Level: logger.TRACE}
			}
		} else {
			// TARGET=LEVEL
			l := strings.TrimSpace(parts[1])

			if lv, err := parseLevel(l); err == nil {
				directive = Directive{Name: t, Level: lv}
			} else {
				directive = Directive{Name: t, Level: logger.TRACE}
			}
		}

		directives = append(directives, directive)
	}

	sort.Sort(sort.Reverse(ByName(directives)))
	return directives
}

func colorize(level logger.Level, s string) string {
	switch level {
	case logger.TRACE:
		return Purple.Paint(s)
	case logger.DEBUG:
		return Blue.Paint(s)
	case logger.INFO:
		return Green.Paint(s)
	case logger.WARN:
		return Yellow.Paint(s)
	case logger.ERROR:
		return Red.Paint(s)
	default:
		return Green.Paint(s)
	}
}
