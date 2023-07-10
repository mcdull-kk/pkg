package log

import (
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	initLoOonce   sync.Once
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
)

func InitLogger(fileName string, level string) {
	initLoOonce.Do(func() {
		initLogger(fileName, level)
	})
}

func initLogger(fileName string, level string) {
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
	logger = zap.New(core).WithOptions(zap.AddCallerSkip(1), zap.AddCaller())
	sugaredLogger = logger.Sugar()
}

func Debug(msg string, fields ...zap.Field)  { logger.Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)   { logger.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)   { logger.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field)  { logger.Error(msg, fields...) }
func DPanic(msg string, fields ...zap.Field) { logger.DPanic(msg, fields...) }
func Panic(msg string, fields ...zap.Field)  { logger.Panic(msg, fields...) }
func Fatal(msg string, fields ...zap.Field)  { logger.Fatal(msg, fields...) }

func Debugw(msg string, keysAndValues ...interface{}) { sugaredLogger.Debugw(msg, keysAndValues...) }
func Infow(msg string, keysAndValues ...interface{})  { sugaredLogger.Infow(msg, keysAndValues...) }
func Warnw(msg string, keysAndValues ...interface{})  { sugaredLogger.Warnw(msg, keysAndValues...) }
func Errorw(msg string, keysAndValues ...interface{}) { sugaredLogger.Errorw(msg, keysAndValues...) }
func Panicw(msg string, keysAndValues ...interface{}) { sugaredLogger.Panicw(msg, keysAndValues...) }
func Fatalw(msg string, keysAndValues ...interface{}) { sugaredLogger.Fatalw(msg, keysAndValues...) }

func Debugf(msg string, keysAndValues ...interface{}) { sugaredLogger.Debugf(msg, keysAndValues...) }
func Infof(msg string, keysAndValues ...interface{})  { sugaredLogger.Infof(msg, keysAndValues...) }
func Warnf(msg string, keysAndValues ...interface{})  { sugaredLogger.Warnf(msg, keysAndValues...) }
func Errorf(msg string, keysAndValues ...interface{}) { sugaredLogger.Errorf(msg, keysAndValues...) }
func Panicf(msg string, keysAndValues ...interface{}) { sugaredLogger.Panicf(msg, keysAndValues...) }
func Fatalf(msg string, keysAndValues ...interface{}) { sugaredLogger.Fatalf(msg, keysAndValues...) }

func Err(keysAndValues ...interface{}) { sugaredLogger.Errorf("error", keysAndValues...) }

func Sync() { logger.Sync() }

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
