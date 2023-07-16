package logrus

import (
	"os"

	"github.com/mcdull-kk/pkg/log"
	"github.com/sirupsen/logrus"
)

var _ log.Logger = (*logrusLogger)(nil)

type logrusLogger struct {
	log *logrus.Logger
}

func NewLogrusLogger(fileName string, level string) log.Logger {
	if level == "" {
		level = "info"
	}
	var logger *logrusLogger
	lv, _ := logrus.ParseLevel(level)
	if fileName != "/dev/stdout" {

	} else {
		logger = &logrusLogger{
			log: &logrus.Logger{
				Level: lv,
				Out:   os.Stdout,
			},
		}
	}
	logger.log.Formatter = &logrus.JSONFormatter{}
	return logger
}

func (l *logrusLogger) Close() error {
	return nil
}

func (l *logrusLogger) Log(level log.Level, keyvals ...any) (err error) {
	var (
		logrusLevel logrus.Level
		fields      logrus.Fields = make(map[string]any)
		msg         string
	)

	switch level {
	case log.DebugLevel:
		logrusLevel = logrus.DebugLevel
	case log.InfoLevel:
		logrusLevel = logrus.InfoLevel
	case log.WarnLevel:
		logrusLevel = logrus.WarnLevel
	case log.ErrorLevel:
		logrusLevel = logrus.ErrorLevel
	case log.FatalLevel:
		logrusLevel = logrus.FatalLevel
	default:
		logrusLevel = logrus.DebugLevel
	}

	if logrusLevel > l.log.Level {
		return
	}

	if len(keyvals) == 0 {
		return nil
	}
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "")
	}
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}
		if key == logrus.FieldKeyMsg {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		fields[key] = keyvals[i+1]
	}

	if len(fields) > 0 {
		l.log.WithFields(fields).Log(logrusLevel, msg)
	} else {
		l.log.Log(logrusLevel, msg)
	}

	return
}
