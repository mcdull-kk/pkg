package etcd

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/mcdull-kk/pkg/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type (
	etcd struct {
		client  *clientv3.Client
		options *options
	}
)

func NewSource(opts ...Option) (config.Source, *clientv3.Client) {
	options := &options{
		Config: clientv3.Config{},
		ctx:    context.Background(),
		path:   "",
		prefix: false,
	}
	for _, opt := range opts {
		opt(options)
	}

	if options.path == "" {
		panic("etcd path invalid")
	}

	client, err := clientv3.New(options.Config)
	if err != nil {
		panic(err)
	}
	return &etcd{client: client, options: options}, client
}

func (e *etcd) Load() ([]*config.KeyValue, error) {
	var opts []clientv3.OpOption
	if e.options.prefix {
		opts = append(opts, clientv3.WithPrefix())
	}
	rsp, err := e.client.Get(e.options.ctx, e.options.path, opts...)
	if err != nil {
		return nil, err
	}
	kvs := make([]*config.KeyValue, 0, len(rsp.Kvs))
	for _, item := range rsp.Kvs {
		k := string(item.Key)
		kvs = append(kvs, &config.KeyValue{
			Key:    k,
			Value:  item.Value,
			Format: strings.TrimPrefix(filepath.Ext(k), "."),
		})
	}
	return kvs, nil
}

func (e *etcd) Watch() (w config.Watcher, err error) {
	return newWatcher(e), nil
}

func (e *etcd) Close() (err error) {
	return e.client.Close()
}
