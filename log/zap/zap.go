package zap

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mcdull-kk/pkg/log"
	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var (
	_              log.Logger = (*Logger)(nil)
	_log           *Logger
	initZapLogOnce sync.Once
)

type Logger struct {
	log *zap.SugaredLogger
}

func (l *Logger) Log(level log.Level, args ...interface{}) error {
	switch level {
	case log.DebugLevel:
		l.log.Debug(args...)
	case log.InfoLevel:
		l.log.Info(args...)
	case log.WarnLevel:
		l.log.Warn(args...)
	case log.ErrorLevel:
		l.log.Error(args...)
	case log.FatalLevel:
		l.log.Fatal(args...)
	}
	return nil
}

func (l *Logger) Close() error {
	return l.log.Sync()
}

func NewLogger(fileName string, level string) *Logger {
	initZapLogOnce.Do(func() {
		initZapLogger(fileName, level)
	})
	return _log
}

func initZapLogger(fileName string, level string) {
	if level == "" {
		level = "info"
	}
	ws := getLogWriter(fileName)
	ecf := zap.NewProductionEncoderConfig()
	ecf.FunctionKey = "func"
	ecf.EncodeTime = zapcore.ISO8601TimeEncoder
	ecf.ConsoleSeparator = " "
	ecf.EncodeLevel = zapcore.LowercaseLevelEncoder
	ecf.EncodeCaller = zapcore.ShortCallerEncoder

	lv, _ := zapcore.ParseLevel(level)
	core := zapcore.NewCore(
		EncodeWrapper{zapcore.NewConsoleEncoder(ecf)},
		&zapcore.BufferedWriteSyncer{
			WS:            ws,
			Size:          0,
			FlushInterval: time.Second * 1,
			Clock:         nil,
		},
		lv,
	)
	_log = &Logger{
		log: zap.New(core).WithOptions(zap.AddCallerSkip(1), zap.AddCaller()).Sugar(),
	}
	log.NewLogger(_log)
}

type EncodeWrapper struct {
	zapcore.Encoder
}

func (e EncodeWrapper) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	function := ent.Caller.Function
	count := 0
	index := strings.IndexFunc(function, func(r rune) bool {
		if r == '/' {
			count++
		}
		if count >= 3 {
			return true
		}
		return false
	})
	function = function[index+1:]

	ent.Caller.Function = function
	return e.Encoder.EncodeEntry(ent, fields)
}

// Save file log cut
func getLogWriter(fileName string) zapcore.WriteSyncer {
	if fileName != "/dev/stdout" {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   fileName, // Log name
			MaxSize:    10,       // File content size, MB
			MaxBackups: 5,        // Maximum number of old files retained
			MaxAge:     30,       // Maximum number of days to keep old files
			Compress:   true,     // Is the file compressed
		}
		return zapcore.AddSync(lumberJackLogger)
	}
	return zapcore.AddSync(os.Stdout)
}
