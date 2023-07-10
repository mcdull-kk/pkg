package redis

import (
	"errors"
	"fmt"
	"time"

	red "github.com/go-redis/redis"
)

const (
	// ClusterType means redis cluster.
	ClusterType = "cluster"
	// NodeType means redis node.
	NodeType      = "node"
	slowThreshold = time.Millisecond * 100
)

var (
	// ErrEmptyHost is an error that indicates no redis host is set.
	ErrEmptyHost = errors.New("empty redis host")
	// ErrEmptyType is an error that indicates no redis type is set.
	ErrEmptyType = errors.New("empty redis type")
)

type (
	// A RedisConf is a redis config.
	RedisConf struct {
		Addr      string `yaml:"addr" json:"addr"`
		Type      string `yaml:"type" json:",default=node,options=node|cluster"`
		Pass      string `yaml:"pass" json:",optional"`
		Tls       bool   `yaml:"tls" json:",default=false,options=true|false"`
		Collector bool   `yaml:"collector" json:",default=false,options=true|false"`
	}

	// Option defines the method to customize a Redis.
	Option func(r *Redis)

	RedisClient interface {
		red.Cmdable
	}

	Redis struct {
		Addr      string
		Type      string
		Pass      string
		tls       bool
		collector bool
	}
)

func (rc RedisConf) NewRedis() *Redis {
	var opts []Option
	if rc.Type == ClusterType {
		opts = append(opts, Cluster())
	}
	if len(rc.Pass) > 0 {
		opts = append(opts, WithPass(rc.Pass))
	}
	if rc.Tls {
		opts = append(opts, WithTLS())
	}
	if rc.Collector {
		opts = append(opts, WithCollector())
	}

	return New(rc.Addr, opts...)
}

// Validate validates the RedisConf.
func (rc RedisConf) Validate() error {
	if len(rc.Addr) == 0 {
		return ErrEmptyHost
	}

	if len(rc.Type) == 0 {
		return ErrEmptyType
	}

	return nil
}

// Cluster customizes the given Redis as a cluster.
func Cluster() Option {
	return func(r *Redis) {
		r.Type = ClusterType
	}
}

// WithPass customizes the given Redis with given password.
func WithPass(pass string) Option {
	return func(r *Redis) {
		r.Pass = pass
	}
}

// WithTLS customizes the given Redis with TLS enabled.
func WithTLS() Option {
	return func(r *Redis) {
		r.tls = true
	}
}

// WithCollector customizes the given Redis with promethus collector enabled.
func WithCollector() Option {
	return func(r *Redis) {
		r.collector = true
	}
}

// New returns a Redis with given options.
func New(addr string, opts ...Option) *Redis {
	r := &Redis{
		Addr: addr,
		Type: NodeType,
	}

	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewRedis returns a Redis.
func NewRedis(redisAddr, redisType string, redisPass ...string) *Redis {
	var opts []Option
	if redisType == ClusterType {
		opts = append(opts, Cluster())
	}
	for _, v := range redisPass {
		opts = append(opts, WithPass(v))
	}

	return New(redisAddr, opts...)
}

// Client Get Redis Client
func (r *Redis) Client() (RedisClient, error) {
	switch r.Type {
	case ClusterType:
		return getCluster(r)
	case NodeType:
		return getClient(r)
	default:
		return nil, fmt.Errorf("redis type '%s' is not supported", r.Type)
	}
}
