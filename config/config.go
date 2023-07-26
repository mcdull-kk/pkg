package config

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mcdull-kk/pkg/codec"
	"github.com/mcdull-kk/pkg/log"
)

var (
	_ Config = (*config)(nil)
	// ErrNotFound is key not found.
	ErrNotFound = errors.New("key not found")
)

type (
	// Observer is config observer.
	Observer func(string, *atomic.Value)

	Config interface {
		Load() error
		Scan(v interface{}) error
		Value(key string) *atomic.Value
		Watch(key string, o Observer) error
		Close() error
	}

	config struct {
		opts      options
		reader    Reader
		cached    sync.Map
		observers sync.Map
		watchers  []Watcher
	}
)

// New a config with options.
func New(opts ...Option) Config {
	o := options{
		decoder:  defaultDecoder,
		resolver: defaultResolver,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &config{
		opts:   o,
		reader: newReader(o),
	}
}

func (c *config) Load() error {
	for _, src := range c.opts.sources {
		kvs, err := src.Load()
		if err != nil {
			return err
		}
		for _, v := range kvs {
			log.Debugf("config loaded: %s format: %s", v.Key, v.Format)
		}
		if err = c.reader.Merge(kvs...); err != nil {
			log.Errorf("failed to merge config source: %v", err)
			return err
		}
		w, err := src.Watch()
		if err != nil {
			log.Errorf("failed to watch config source: %v", err)
			return err
		}
		c.watchers = append(c.watchers, w)
		go c.watch(w)
	}
	if err := c.reader.Resolve(); err != nil {
		log.Errorf("failed to resolve config source: %v", err)
		return err
	}
	return nil
}

func (c *config) Value(key string) *atomic.Value {
	if v, ok := c.cached.Load(key); ok {
		return v.(*atomic.Value)
	}
	if v, ok := c.reader.Value(key); ok {
		c.cached.Store(key, v)
		return v
	}
	return nil
}

func (c *config) Scan(v any) error {
	data, err := c.reader.Source()
	if err != nil {
		return err
	}

	return codec.GetCodec(codec.JsonName).Unmarshal(data, v)
}

func (c *config) Watch(key string, o Observer) error {
	if v := c.Value(key); v.Load() == nil {
		return ErrNotFound
	}
	c.observers.Store(key, o)
	return nil
}

func (c *config) Close() error {
	for _, w := range c.watchers {
		if err := w.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (c *config) watch(w Watcher) {
	for {
		kvs, err := w.Next()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Infof("watcher's ctx cancel : %v", err)
				return
			}
			time.Sleep(time.Second)
			log.Errorf("failed to watch next config: %v", err)
			continue
		}
		if err := c.reader.Merge(kvs...); err != nil {
			log.Errorf("failed to merge next config: %v", err)
			continue
		}
		if err := c.reader.Resolve(); err != nil {
			log.Errorf("failed to resolve next config: %v", err)
			continue
		}
		c.cached.Range(func(key, value interface{}) bool {
			k := key.(string)
			v := value.(*atomic.Value)
			if n, ok := c.reader.Value(k); ok && reflect.TypeOf(n.Load()) == reflect.TypeOf(v.Load()) && !reflect.DeepEqual(n.Load(), v.Load()) {
				v.Store(n.Load())
				if o, ok := c.observers.Load(k); ok {
					o.(Observer)(k, v)
				}
			}
			return true
		})
	}
}
