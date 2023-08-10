package etcd

import (
	"context"
	"crypto/tls"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type (
	Option func(*options)

	options struct {
		clientv3.Config
		ctx    context.Context
		path   string
		prefix bool
	}
)

func WithEndpoints(endpoints []string) Option {
	return func(o *options) {
		o.Endpoints = endpoints
	}
}

func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(o *options) {
		o.DialTimeout = dialTimeout
	}
}

func WithTLS(tls *tls.Config) Option {
	return func(o *options) {
		o.TLS = tls
	}
}

func WithDialOptions(dialOptions []grpc.DialOption) Option {
	return func(o *options) {
		o.DialOptions = dialOptions
	}
}

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
