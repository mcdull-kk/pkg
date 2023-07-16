package zap

import (
	slog "log"
	"testing"

	"github.com/mcdull-kk/pkg/log"
)

func TestLogWithFilter(t *testing.T) {
	logger := NewZapLogger("/dev/stdout", "debug")
	log.SetLogger(log.NewFilter(logger,
		log.FilterLevel(log.InfoLevel),
		log.FilterKey("password"),
		log.FilterValue("haha"),
		log.FilterFunc(testFilterFunc)))
	defer func() {
		logger.Close()
	}()

	log.Debug("hello")
	log.Infow("password", "12345")
	log.Warn("werq")
	log.Infow("phone", "123456")
	log.Info("sdfafdafaff")
	log.Infof("比好 %s", "哈哈")
}

func testFilterFunc(level log.Level, keyvals ...any) bool {
	if level == log.WarnLevel {
		return true
	}
	for i := 0; i < len(keyvals); i++ {
		if keyvals[i] == "phone" {
			keyvals[i+1] = "***"
		}
	}
	return false
}

func TestStdLog(t *testing.T) {
	log.SetLogger(log.NewFilter(log.NewStdLogger(slog.Writer()),
		log.FilterLevel(log.InfoLevel),
		log.FilterKey("password"),
		log.FilterValue("haha"),
		log.FilterFunc(testFilterFunc)))

	log.Debug("hello")
	log.Infow("password", "12345")
	log.Warn("werq")
	log.Infow("phone", "123456")
	log.Info("sdfafdafaff")
	log.Infof("比好 %s", "哈哈")
}
