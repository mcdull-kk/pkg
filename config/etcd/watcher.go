package etcd

import (
	"context"

	"github.com/mcdull-kk/pkg/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var _ config.Watcher = (*watcher)(nil)

type (
	watcher struct {
		etcd   *etcd
		ch     clientv3.WatchChan
		ctx    context.Context
		cancel context.CancelFunc
	}
)

func newWatcher(e *etcd) config.Watcher {
	ctx, cancel := context.WithCancel(context.Background())
	w := &watcher{
		etcd:   e,
		ctx:    ctx,
		cancel: cancel,
	}

	var opts []clientv3.OpOption
	if e.options.prefix {
		opts = append(opts, clientv3.WithPrefix())
	}
	w.ch = e.client.Watch(e.options.ctx, e.options.path, opts...)
	return w
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
	select {
	case resp := <-w.ch:
		if resp.Err() != nil {
			return nil, resp.Err()
		}
		return w.etcd.Load()
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	}
}

func (w *watcher) Stop() error {
	w.cancel()
	return nil
}
