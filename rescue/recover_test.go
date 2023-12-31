package rescue

import (
	"sync/atomic"
	"testing"

	"github.com/mcdull-kk/pkg/log/zap"
	"github.com/stretchr/testify/assert"
)

func TestRescue(t *testing.T) {
	zap.NewZapLogger("/dev/stdout", "debug")
	var count int32
	assert.NotPanics(t, func() {
		defer Recover(func() {
			atomic.AddInt32(&count, 2)
		}, func() {
			atomic.AddInt32(&count, 3)
		})

		panic("hello")
	})
	assert.Equal(t, int32(5), atomic.LoadInt32(&count))
}
