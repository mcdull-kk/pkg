package consul

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/consul/api"
)

type (
	Option func(o *options)

	options struct {
		*api.Config
		ctx  context.Context
		path string
	}
)

func WithAddr(addr string) Option {
	return func(o *options) {
		o.Address = addr
	}
}

func WithScheme(scheme string) Option {
	return func(o *options) {
		o.Scheme = scheme
	}
}

func WithDatacenter(datacenter string) Option {
	return func(o *options) {
		o.Datacenter = datacenter
	}
}

func WithTransport(t *http.Transport) Option {
	return func(o *options) {
		o.Transport = t
	}
}

func WithHttpClient(c *http.Client) Option {
	return func(o *options) {
		o.HttpClient = c
	}
}

func WithHttpAuth(ha *api.HttpBasicAuth) Option {
	return func(o *options) {
		o.HttpAuth = ha
	}
}

func WithWaitTime(wt time.Duration) Option {
	return func(o *options) {
		o.WaitTime = wt
	}
}

func WithToken(token string) Option {
	return func(o *options) {
		o.Token = token
	}
}

func WithTokenFile(tf string) Option {
	return func(o *options) {
		o.TokenFile = tf
	}
}

func WithTLSConfig(tls api.TLSConfig) Option {
	return func(o *options) {
		o.TLSConfig = tls
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
