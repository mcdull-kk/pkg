package apollo

import (
	"context"

	"github.com/mcdull-kk/pkg/config"
)

type (
	watcher struct {
		out <-chan []*config.KeyValue

		ctx      context.Context
		cancelFn func()
	}

	changeListener struct {
		in     chan<- []*config.KeyValue
		apollo *apollo
	}
)

func newWatcher(a *apollo)
