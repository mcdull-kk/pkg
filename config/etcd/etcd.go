package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type (
	etcd struct {
		client  clientv3.Client
		options *options
	}

	Option func(*options)

	options struct {
		*clientv3.Config
		ctx    context.Context
		path   string
		prefix bool
	}
)

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func WithPath(p string) Option {
	return func(o *options) {
		o.path = p
	}
}

func WithPrefix(prefix bool) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}
