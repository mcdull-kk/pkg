package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
)

// Level is a logger level.
type Level int8

const (
	// DebugLevel is logger debug level.
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return ""
	}
}

// ParseLevel parses a level string into a logger Level value.
func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return DebugLevel
	case "INFO":
		return InfoLevel
	case "WARN":
		return WarnLevel
	case "ERROR":
		return ErrorLevel
	case "FATAL":
		return FatalLevel
	}
	return InfoLevel
}

// Logger is a logger interface.
type Logger interface {
	Log(level Level, keyvals ...any) error
}

type (
	stdLogger struct {
		log  *log.Logger
		pool *sync.Pool
	}

	Filter struct {
		logger Logger
		level  Level
		key    map[any]struct{}
		value  map[any]struct{}
		filter func(level Level, keyvals ...any) bool
	}
	FilterOption func(*Filter)
)

func NewStdLogger(w io.Writer) Logger {
	return &stdLogger{
		log: log.New(w, "", 0),
		pool: &sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

// Log print the kv pairs log.
func (l *stdLogger) Log(level Level, keyvals ...any) error {
	if len(keyvals) == 0 {
		return nil
	}
	if (len(keyvals) & 1) == 1 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}
	buf := l.pool.Get().(*bytes.Buffer)
	buf.WriteString(level.String())
	for i := 0; i < len(keyvals); i += 2 {
		_, _ = fmt.Fprintf(buf, " %s=%v", keyvals[i], keyvals[i+1])
	}
	_ = l.log.Output(4, buf.String()) //nolint:gomnd
	buf.Reset()
	l.pool.Put(buf)
	return nil
}

func (l *stdLogger) Close() error {
	return nil
}

// FilterLevel with filter level.
func FilterLevel(level Level) FilterOption {
	return func(opts *Filter) {
		opts.level = level
	}
}

// FilterKey with filter key.
func FilterKey(key ...string) FilterOption {
	return func(o *Filter) {
		for _, v := range key {
			o.key[v] = struct{}{}
		}
	}
}

// FilterValue with filter value.
func FilterValue(value ...string) FilterOption {
	return func(o *Filter) {
		for _, v := range value {
			o.value[v] = struct{}{}
		}
	}
}

// FilterFunc with filter func.
func FilterFunc(f func(level Level, keyvals ...any) bool) FilterOption {
	return func(o *Filter) {
		o.filter = f
	}
}

func NewFilter(logger Logger, opts ...FilterOption) Logger {
	f := &Filter{
		logger: logger,
		key:    make(map[any]struct{}),
		value:  make(map[any]struct{}),
	}
	for _, o := range opts {
		o(f)
	}
	return f
}

func (f *Filter) Log(level Level, keyvals ...any) error {
	if level < f.level {
		return nil
	}
	if f.filter != nil && f.filter(level, keyvals...) {
		return nil
	}
	if len(f.key) > 0 || len(f.value) > 0 {
		for i := 0; i < len(keyvals); i += 2 {
			v := i + 1
			if v >= len(keyvals) {
				continue
			}
			if _, ok := f.key[keyvals[i]]; ok {
				keyvals[v] = "***"
			}
			if _, ok := f.value[keyvals[v]]; ok {
				keyvals[v] = "***"
			}
		}
	}
	return f.logger.Log(level, keyvals...)
}
