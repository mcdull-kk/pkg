package log

import (
	"testing"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	InitLogger("/dev/stdout", "debug")
	defer Sync()

	Debug("info", zap.String("app", "start ok"), zap.Int("ik", 1))
}
