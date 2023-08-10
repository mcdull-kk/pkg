package consul

import (
	"context"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/mcdull-kk/pkg/config"
	"github.com/mcdull-kk/pkg/rescue"
)

var _ config.Watcher = (*watcher)(nil)

type watcher struct {
	consul *consul
	ch     chan interface{}
	wp     *watch.Plan
	ctx    context.Context
	cancel context.CancelFunc
}

func newWatcher(c *consul) (*watcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	w := &watcher{
		consul: c,
		ch:     make(chan interface{}),
		ctx:    ctx,
		cancel: cancel,
	}
	wp, err := watch.Parse(map[string]interface{}{"type": "keyprefix", "prefix": c.options.path})
	if err != nil {
		return nil, err
	}
	wp.Handler = w.handle
	w.wp = wp

	rescue.GoSafe(func() {
		err := wp.RunWithClientAndHclog(c.client, nil)
		if err != nil {
			panic(err)
		}
	})
	return w, nil
}

func (w *watcher) handle(_ uint64, data interface{}) {
	if data == nil {
		return
	}

	_, ok := data.(api.KVPairs)
	if !ok {
		return
	}

	w.ch <- struct{}{}
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
	select {
	case <-w.ch:
		return w.consul.Load()
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	}
}

func (w *watcher) Stop() error {
	w.wp.Stop()
	w.cancel()
	return nil
}
