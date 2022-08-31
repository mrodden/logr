package fmt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/mattn/go-isatty"

	"github.com/mrodden/logr/logger"
)

const (
	ISO = "2006-01-02T15:04:05.999Z"
)

func GoRoutineID() string {
	return string(bytes.Fields(debug.Stack())[1])
}

type SyncWriter struct {
	w    io.Writer
	mu   sync.Mutex
	ansi bool
}

func NewSyncWriter(w io.Writer) *SyncWriter {
	return &SyncWriter{w: w, ansi: SupportsAnsi(w)}
}

func (w *SyncWriter) HasAnsiSupport() bool {
	return w.ansi
}

func (w *SyncWriter) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.w.Write(b)
}

func SupportsAnsi(w io.Writer) bool {
	f, ok := w.(*os.File)
	if runtime.GOOS == "windows" {
		ok = false
	}
	return ok && isatty.IsTerminal(f.Fd())
}

type FormatEvent = func(Format, *SyncWriter, logger.Record) error

type Format struct {
	Format             FormatEvent
	DisplayGoroutineId bool
}

func (f Format) FormatEvent(w *SyncWriter, event logger.Record) error {
	return f.Format(f, w, event)
}

func Full(f Format, w *SyncWriter, event logger.Record) error {
	level := event.Metadata().Level()
	comp := event.Metadata().Target()

	dimmed := Style{}
	if w.HasAnsiSupport() {
		dimmed = Style{}.Dimmed()
	}

	buf := make([]byte, 0)

	ts := dimmed.Paint(fmt.Sprintf("%-24s ", time.Now().UTC().Format(ISO)))
	buf = append(buf, ts...)

	if f.DisplayGoroutineId {
		thr := dimmed.Paint(fmt.Sprintf("%02v ", GoRoutineID()))
		buf = append(buf, thr...)
	}

	cond := fmt.Sprintf("%5s ", logger.Ltoa(level))
	if w.HasAnsiSupport() {
		cond = colorize(level, cond)
	}
	buf = append(buf, cond...)

	comp = dimmed.Paint(comp)
	buf = append(buf, comp...)

	buf = append(buf, dimmed.Paint(":")...)
	buf = append(buf, ' ')

	s := fmt.Sprint(event.Args()...)
	buf = append(buf, s...)

	if len(s) == 0 || s[len(s)-1] != '\n' {
		buf = append(buf, '\n')
	}

	_, err := w.Write(buf)
	return err
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
