package consul

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/mcdull-kk/pkg/config"
)

type (
	consul struct {
		client  *api.Client
		options *options
	}
)

func NewSource(opts ...Option) config.Source {
	options := &options{
		Config: &api.Config{},
		ctx:    context.Background(),
		path:   "",
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Address == "" {
		panic("consul address is empty")
	}
	if options.path == "" {
		panic("consul path is empty")
	}

	client, err := api.NewClient(options.Config)
	if err != nil {
		panic(err)
	}
	return &consul{client: client, options: options}
}

func (c *consul) Load() (kv []*config.KeyValue, err error) {
	kvs, _, err := c.client.KV().List(c.options.path, nil)
	if err != nil {
		return nil, err
	}
	pathPrefix := c.options.path
	if !strings.HasPrefix(c.options.path, "/") {
		pathPrefix += "/"
	}
	kv = make([]*config.KeyValue, 0)
	for _, item := range kvs {
		k := strings.TrimPrefix(item.Key, pathPrefix)
		if k == "" {
			continue
		}
		kv = append(kv, &config.KeyValue{
			Key:    k,
			Value:  item.Value,
			Format: strings.TrimPrefix(filepath.Ext(k), "."),
		})
	}
	return
}

func (c *consul) Watch() (config.Watcher, error) {
	return newWatcher(c)
}

func (c *consul) Close() (err error) {
	return nil
}
