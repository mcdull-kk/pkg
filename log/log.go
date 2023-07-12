package log

import (
	"fmt"
	"strings"
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
	Log(level Level, keyvals ...interface{}) error
	Close() error
}

// Option is log option.
type Option func(*logger)

// WithSprint with sprint
func WithSprint(sprint func(...interface{}) string) Option {
	return func(l *logger) {
		l.sprint = sprint
	}
}

// WithSprintf with sprintf
func WithSprintf(sprintf func(format string, a ...interface{}) string) Option {
	return func(l *logger) {
		l.sprintf = sprintf
	}
}

var (
	l *logger
)

type logger struct {
	Logger
	sprint  func(...interface{}) string
	sprintf func(format string, a ...interface{}) string
}

func NewLogger(log Logger, opts ...Option) {
	l = &logger{
		Logger:  log,
		sprint:  fmt.Sprint,
		sprintf: fmt.Sprintf,
	}
	for _, o := range opts {
		o(l)
	}
}

func Debug(a ...interface{})                 { _ = l.Log(DebugLevel, l.sprint(a...)) }
func Debugf(format string, a ...interface{}) { _ = l.Log(DebugLevel, l.sprintf(format, a...)) }
func Debugw(keyvals ...interface{})          { _ = l.Log(DebugLevel, keyvals...) }
func Info(a ...interface{})                  { _ = l.Log(InfoLevel, l.sprint(a...)) }
func Infof(format string, a ...interface{})  { _ = l.Log(InfoLevel, l.sprintf(format, a...)) }
func Infow(keyvals ...interface{})           { _ = l.Log(InfoLevel, keyvals...) }
func Warn(a ...interface{})                  { _ = l.Log(WarnLevel, l.sprint(a...)) }
func Warnf(format string, a ...interface{})  { _ = l.Log(WarnLevel, l.sprintf(format, a...)) }
func Warnw(keyvals ...interface{})           { _ = l.Log(WarnLevel, keyvals...) }
func Error(a ...interface{})                 { _ = l.Log(ErrorLevel, l.sprint(a...)) }
func Errorf(format string, a ...interface{}) { _ = l.Log(ErrorLevel, l.sprintf(format, a...)) }
func Errorw(keyvals ...interface{})          { _ = l.Log(ErrorLevel, keyvals...) }
func Fatal(a ...interface{})                 { _ = l.Log(FatalLevel, l.sprint(a...)) }
func Fatalf(format string, a ...interface{}) { _ = l.Log(FatalLevel, l.sprintf(format, a...)) }
func Fatalw(keyvals ...interface{})          { _ = l.Log(FatalLevel, keyvals...) }
func Close()                                 { _ = l.Close() }
