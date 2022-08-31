package logger

import (
	"errors"
	"sync"
)

type Level int

const (
	ERROR Level = iota
	WARN
	INFO
	DEBUG
	TRACE
)

func (l Level) String() string {
	return Ltoa(l)
}

func Ltoa(level Level) string {
	switch level {
	case ERROR:
		return "ERROR"
	case WARN:
		return "WARN"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	default:
		return "UNDEF"
	}
}

type Logger interface {
	Enabled(Metadata) bool
	Log(Record)
	Flush()
}

type Record struct {
	metadata    Metadata
	args        []any
	packagePath string
	file        string
	line        uint32
}

func (r *Record) Metadata() *Metadata { return &r.metadata }
func (r *Record) Args() []any         { return r.args }
func (r *Record) Filename() string    { return r.file }
func (r *Record) LineNumber() uint32  { return r.line }

// RecordBuilder implements the builder pattern for Record objects
// For example:
//
//	rec := logger.NewRecordBuilder().
//		Level(level).
//		Target(pkg).
//		Args(s).
//		PackagePath(pkg).
//		File(file).
//		Line(line).
//		Build()
//
type RecordBuilder struct {
	record Record
}

func NewRecordBuilder() *RecordBuilder {
	return &RecordBuilder{}
}

func (rb *RecordBuilder) Level(level Level) *RecordBuilder {
	rb.record.metadata.level = level
	return rb
}

func (rb *RecordBuilder) Target(target string) *RecordBuilder {
	rb.record.metadata.target = target
	return rb
}

func (rb *RecordBuilder) Args(v ...any) *RecordBuilder {
	rb.record.args = v
	return rb
}

func (rb *RecordBuilder) PackagePath(path string) *RecordBuilder {
	rb.record.packagePath = path
	return rb
}

func (rb *RecordBuilder) File(file string) *RecordBuilder {
	rb.record.file = file
	return rb
}

func (rb *RecordBuilder) Line(line uint) *RecordBuilder {
	rb.record.line = uint32(line)
	return rb
}

func (rb *RecordBuilder) Build() Record {
	return rb.record
}

type Metadata struct {
	level  Level
	target string
}

func (m *Metadata) Level() Level   { return m.level }
func (m *Metadata) Target() string { return m.target }

var mu = sync.Mutex{}
var global Logger = nil

func Global() Logger {
	return global
}

func SetDefaultLogger(l Logger) error {
	mu.Lock()
	defer mu.Unlock()

	if global == nil {
		global = l
		return nil
	} else {
		return errors.New("global logger already set")
	}
}
