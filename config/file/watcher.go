package file

import (
	"context"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/mcdull-kk/pkg/config"
)

var _ config.Watcher = (*watcher)(nil)

type watcher struct {
	f  *file
	fw *fsnotify.Watcher

	ctx    context.Context
	cancel context.CancelFunc
}

func newWatcher(f *file) (config.Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := fw.Add(f.path); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &watcher{f: f, fw: fw, ctx: ctx, cancel: cancel}, nil
}

func (w *watcher) Next() (kvs []*config.KeyValue, err error) {
	select {
	case <-w.ctx.Done():
		err = w.ctx.Err()
		return
	case err = <-w.fw.Errors:
		return
	case event := <-w.fw.Events:
		if event.Op == fsnotify.Rename {
			if _, err := os.Stat(event.Name); err == nil || os.IsExist(err) {
				if err := w.fw.Add(event.Name); err != nil {
					return nil, err
				}
			}
		}
		fi, err := os.Stat(w.f.path)
		if err != nil {
			return nil, err
		}
		path := w.f.path
		if fi.IsDir() {
			path = filepath.Join(w.f.path, filepath.Base(event.Name))
		}
		kv, err := w.f.loadFile(path)
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, kv)
		return kvs, nil
	}
}

func (w *watcher) Stop() error {
	w.cancel()
	return w.fw.Close()
}
