package logrus

import (
	"testing"

	"github.com/mcdull-kk/pkg/log"
)

func TestLog(t *testing.T) {
	NewLogger("/dev/stdout", "debug")

	log.Debug("adsf")
}
