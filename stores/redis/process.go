package redis

import (
	"fmt"
	"strings"
	"time"

	red "github.com/go-redis/redis"
	"github.com/mcdull-kk/pkg/codec"
	"github.com/mcdull-kk/pkg/log"
	"go.uber.org/zap"
)

func process(proc func(red.Cmder) error) func(red.Cmder) error {
	return func(cmd red.Cmder) error {
		st := time.Now()
		defer func() {
			duration := time.Since(st)
			if duration > slowThreshold {
				var buf strings.Builder
				for i, arg := range cmd.Args() {
					if i > 0 {
						buf.WriteByte(' ')
					}
					buf.WriteString(codec.Repr(arg))
				}
				log.Warn("[REDIS]", zap.String("slowcall on executing", buf.String()), zap.String("costtime", fmt.Sprintf("%.1fms", float32(duration)/float32(time.Millisecond))))
			}
		}()

		return proc(cmd)
	}
}
