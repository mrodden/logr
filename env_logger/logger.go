package env_logger

import (
	"fmt"
	"os"
	"sort"
	"strings"

	format "github.com/mrodden/logr/env_logger/fmt"
	"github.com/mrodden/logr/logger"
)

type builder struct {
	logger envLogger
}

// Builder creates a new builder for configuring a logger.
func Builder() builder {
	return builder{
		logger: envLogger{
			directives: DefaultEnvFilter(),
			format:     DefaultFormat(),
			writer:     format.NewSyncWriter(os.Stderr),
		},
	}
}

// JSON sets the formatting of the logger to JSON output.
func (b builder) JSON() builder {
	b.logger.format.Format = format.JSON
	b.logger.format.DisplayGoroutineId = false
	return b
}

// WithGoroutineID sets the logging format to include the goroutine ID of the caller
func (b builder) WithGoroutineId() builder {
	b.logger.format.DisplayGoroutineId = true
	return b
}

// WithEnvFilter can be used to set a custom set of filter Directives
// for the logger created by the builder.
// See EnvFilter and DefaultEnvFilter for building Directives.
func (b builder) WithEnvFilter(d []Directive) builder {
	b.logger.directives = d
	return b
}

// Build finalizes the logger and returns it.
func (b builder) Build() *envLogger {
	return &b.logger
}

// TryInit finalizes the logger and sets it as the default logger.
// This function will return an error if setting the default logger fails.
func (b builder) TryInit() error {
	l := b.Build()

	// set global logr Logger
	return logger.SetDefaultLogger(l)
}

// Init finalizes the logger and sets it as the default logger.
// This function will panic if setting the default logger fails.
func (b builder) Init() {
	err := b.TryInit()
	if err != nil {
		panic(fmt.Sprintf("error during logger init: %v", err))
	}
}

func TryInit() error {
	return Builder().TryInit()
}

func Init() {
	Builder().Init()
}

func DefaultFormat() format.Format {
	fs := strings.ToLower(os.Getenv("GO_LOG_FMT"))
	switch fs {
	case "json":
		return format.Format{Format: format.JSON, DisplayGoroutineId: false}
	default:
		return format.Format{Format: format.Full, DisplayGoroutineId: true}
	}
}

func DefaultEnvFilter() []Directive {
	return EnvFilter(os.Getenv("GO_LOG"))
}

func EnvFilter(env string) []Directive {
	directives := parse(env)

	if len(directives) == 0 {
		directives = append(directives, Directive{Name: "", Level: logger.ERROR})
	}
	return directives
}

type envLogger struct {
	directives []Directive
	format     format.Format
	writer     *format.SyncWriter
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
		_ = l.format.FormatEvent(l.writer, record)
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
