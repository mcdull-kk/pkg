package apollo

import (
	"context"

	"github.com/apolloconfig/agollo/v4/constant"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/mcdull-kk/pkg/codec"
	"github.com/mcdull-kk/pkg/config"
	"github.com/mcdull-kk/pkg/log"
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

func newWatcher(a *apollo) config.Watcher {
	changeCh := make(chan []*config.KeyValue)
	listener := &changeListener{in: changeCh, apollo: a}
	a.client.AddChangeListener(listener)

	ctx, cancel := context.WithCancel(context.Background())
	return &watcher{
		out: changeCh,
		ctx: ctx,
		cancelFn: func() {
			a.client.RemoveChangeListener(listener)
			cancel()
		},
	}
}

func (l *changeListener) OnNewestChange(_ *storage.FullChangeEvent) {}

func (l *changeListener) OnChange(event *storage.ChangeEvent) {
	kv := make([]*config.KeyValue, 0, 2)
	ns := event.Namespace
	fm := configFileformat(ns)

	if fm == constant.JSON || fm == constant.YML || fm == constant.YAML || fm == constant.XML {
		value, err := l.apollo.client.GetConfigCache(ns).Get("content")
		if err != nil {
			log.Warnw("apollo get config failed", "err", err)
			return
		}
		kv = append(kv, &config.KeyValue{
			Key:    ns,
			Value:  []byte(value.(string)),
			Format: format(ns),
		})
	} else {
		next := make(map[string]any)
		for key, change := range event.Changes {
			resolve(genKey(ns, key), change.NewValue, next)
		}
		f := format(ns)
		code := codec.GetCodec(f)
		val, err := code.Marshal(next)
		if err != nil {
			log.Warnf("apollo could not handle namespace %s: %v", ns, err)
			return
		}
		kv = append(kv, &config.KeyValue{
			Key:    ns,
			Value:  val,
			Format: f,
		})
	}

	l.in <- kv
}

// Next will be blocked until the Stop method is called
func (w *watcher) Next() ([]*config.KeyValue, error) {
	select {
	case kv := <-w.out:
		return kv, nil
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	}
}

func (w *watcher) Stop() error {
	if w.cancelFn != nil {
		w.cancelFn()
	}
	return nil
}
