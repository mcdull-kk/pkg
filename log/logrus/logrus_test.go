package logrus

import (
	"testing"

	"github.com/mcdull-kk/pkg/log"
)

func TestLog(t *testing.T) {
	l := NewLogrusLogger("/dev/stdout", "debug")
	log.SetLogger(l)
	log.Debug("adsf")
}
