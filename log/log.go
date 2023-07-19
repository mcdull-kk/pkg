package log

import (
	"fmt"
	"log"
	"sync"
)

var (
	global = &logger{}
)

func init() {
	global.SetLogger(NewStdLogger(log.Writer()))
}

type logger struct {
	lock sync.Mutex
	Logger
}

func (l *logger) SetLogger(logger Logger) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.Logger = logger
}

func SetLogger(logger Logger) {
	global.SetLogger(logger)
}

func GetLogger() Logger {
	return global
}

func GetGlobalLogger() *logger {
	return global
}

func Debug(a ...any) { _ = global.Log(DebugLevel, fmt.Sprint(a...)) }
func Debugf(format string, a ...any) {
	_ = global.Log(DebugLevel, fmt.Sprintf(format, a...))
}
func Debugw(keyvals ...any)         { _ = global.Log(DebugLevel, keyvals...) }
func Info(a ...any)                 { _ = global.Log(InfoLevel, fmt.Sprint(a...)) }
func Infof(format string, a ...any) { _ = global.Log(InfoLevel, fmt.Sprintf(format, a...)) }
func Infow(keyvals ...any)          { _ = global.Log(InfoLevel, keyvals...) }
func Warn(a ...any)                 { _ = global.Log(WarnLevel, fmt.Sprint(a...)) }
func Warnf(format string, a ...any) { _ = global.Log(WarnLevel, fmt.Sprintf(format, a...)) }
func Warnw(keyvals ...any)          { _ = global.Log(WarnLevel, keyvals...) }
func Error(a ...any)                { _ = global.Log(ErrorLevel, fmt.Sprint(a...)) }
func Errorf(format string, a ...any) {
	_ = global.Log(ErrorLevel, fmt.Sprintf(format, a...))
}
func Errorw(keyvals ...any) { _ = global.Log(ErrorLevel, keyvals...) }
func Fatal(a ...any)        { _ = global.Log(FatalLevel, fmt.Sprint(a...)) }
func Fatalf(format string, a ...any) {
	_ = global.Log(FatalLevel, fmt.Sprintf(format, a...))
}
func Fatalw(keyvals ...any) { _ = global.Log(FatalLevel, keyvals...) }

func (l *logger) Debug(a ...any)                 { Debug(a...) }
func (l *logger) Debugf(format string, a ...any) { Debugf(format, a...) }
func (l *logger) Info(keyvals ...any)            { Debugw(keyvals...) }
func (l *logger) Infof(format string, a ...any)  { Infof(format, a...) }
func (l *logger) Infow(keyvals ...any)           { Infow(keyvals...) }
func (l *logger) Warn(a ...any)                  { Warn(a...) }
func (l *logger) Warnf(format string, a ...any)  { Warnf(format, a...) }
func (l *logger) Warnw(keyvals ...any)           { Warnw(keyvals...) }
func (l *logger) Error(a ...any)                 { Error(a...) }
func (l *logger) Errorf(format string, a ...any) { Errorf(format, a...) }
func (l *logger) Errorw(keyvals ...any)          { Errorw(keyvals...) }
func (l *logger) Fatal(a ...any)                 { Fatal(a...) }
func (l *logger) Fatalf(format string, a ...any) { Fatalf(format, a...) }
func (l *logger) Fatalw(keyvals ...any)          { Fatalw(keyvals...) }
