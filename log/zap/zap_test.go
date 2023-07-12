package zap

import (
	"testing"

	"github.com/mcdull-kk/pkg/log"
)

func TestLog(t *testing.T) {
	l := NewLogger("/dev/stdout", "debug")
	defer func() {
		l.Close()
	}()

	log.Debugf("")
}
