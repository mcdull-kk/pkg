package apollo

import (
	"strings"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/constant"
	apolloconfig "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/extension"

	"github.com/mcdull-kk/pkg/codec"
	"github.com/mcdull-kk/pkg/config"
	"github.com/mcdull-kk/pkg/log"
)

type (
	apollo struct {
		client agollo.Client
		opt    *options
	}

	Option func(*apolloconfig.AppConfig)

	options struct {
		*apolloconfig.AppConfig
		originConfig bool
	}

	extParser struct{}
)

func WithAppID(appID string) Option {
	return func(o *apolloconfig.AppConfig) {
		o.AppID = appID
	}
}

func WithIP(ip string) Option {
	return func(o *apolloconfig.AppConfig) {
		o.IP = ip
	}
}

func WithCluster(cluster string) Option {
	return func(o *apolloconfig.AppConfig) {
		o.Cluster = cluster
	}
}

func WithSecret(secret string) Option {
	return func(o *apolloconfig.AppConfig) {
		o.Secret = secret
	}
}

func WithNamespace(namespace string) Option {
	return func(o *apolloconfig.AppConfig) {
		o.NamespaceName = namespace
	}
}

func (parser extParser) Parse(configContent interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{"content": configContent}, nil
}

func NewSource(originConfig bool, opts ...Option) config.Source {
	op := &apolloconfig.AppConfig{}
	for _, o := range opts {
		o(op)
	}

	agollo.SetLogger(log.GetGlobalLogger())

	client, err := agollo.StartWithConfig(func() (*apolloconfig.AppConfig, error) {
		return op, nil
	})
	if err != nil {
		panic(err)
	}

	opt := &options{op, originConfig}
	if opt.originConfig {
		extension.AddFormatParser(constant.JSON, &extParser{})
		extension.AddFormatParser(constant.Properties, &extParser{})
		extension.AddFormatParser(constant.XML, &extParser{})
		extension.AddFormatParser(constant.YML, &extParser{})
		extension.AddFormatParser(constant.YAML, &extParser{})
	}
	return &apollo{client: client, opt: opt}
}

func (e *apollo) Load() (kv []*config.KeyValue, err error) {
	kvs := make([]*config.KeyValue, 0)
	namespaces := strings.Split(e.opt.NamespaceName, ",")

	for _, ns := range namespaces {
		if !e.opt.originConfig {
			kv, err := e.getConfig(ns)
			if err != nil {
				log.Errorf("apollo get config failed，err:%v", err)
				continue
			}
			kvs = append(kvs, kv)
			continue
		}

		fm := configFileformat(ns)
		if fm == constant.JSON || fm == constant.YML || fm == constant.YAML || fm == constant.XML {
			kv, err := e.getOriginConfig(ns)
			if err != nil {
				log.Errorf("apollo get config failed，err:%v", err)
				continue
			}
			kvs = append(kvs, kv)
			continue
		}
		kv, err := e.getConfig(ns)
		if err != nil {
			log.Errorf("apollo get config failed，err:%v", err)
			continue
		}
		kvs = append(kvs, kv)
	}
	return kvs, nil
}

func (e *apollo) Watch() (config.Watcher, error) {
	return newWatcher(e), nil
}

func (e *apollo) getConfig(ns string) (*config.KeyValue, error) {
	next := map[string]interface{}{}
	e.client.GetConfigCache(ns).Range(func(key, value interface{}) bool {
		// all values are out properties format
		resolve(genKey(ns, key.(string)), value, next)
		return true
	})
	f := format(ns)
	code := codec.GetCodec(f)
	val, err := code.Marshal(next)
	if err != nil {
		return nil, err
	}
	return &config.KeyValue{
		Key:    ns,
		Value:  val,
		Format: f,
	}, nil
}

func (e apollo) getOriginConfig(ns string) (*config.KeyValue, error) {
	value, err := e.client.GetConfigCache(ns).Get("content")
	if err != nil {
		return nil, err
	}
	// serialize the namespace content KeyValue into bytes.
	return &config.KeyValue{
		Key:    ns,
		Value:  []byte(value.(string)),
		Format: format(ns),
	}, nil
}

func configFileformat(ns string) constant.ConfigFileFormat {
	arr := strings.Split(ns, ".")
	if len(arr) <= 1 {
		return constant.DEFAULT
	}
	return constant.ConfigFileFormat("." + arr[len(arr)-1])
}

func format(ns string) string {
	arr := strings.Split(ns, ".")
	suffix := arr[len(arr)-1]
	if len(arr) <= 1 || suffix == "properties" {
		return "json"
	}
	fm := constant.ConfigFileFormat("." + suffix)
	if fm != constant.JSON && fm != constant.YAML && fm != constant.XML && fm != constant.YML {
		// fallback
		return "json"
	}
	return suffix
}

// resolve convert kv pair into one map[string]any by split key into different
// map level. such as: app.name = "application" => map[app][name] = "application"
func resolve(key string, value any, target map[string]any) {
	// expand key "aaa.bbb" into map[aaa]map[bbb]any
	keys := strings.Split(key, ".")
	last := len(keys) - 1
	cursor := target

	for i, k := range keys {
		if i == last {
			cursor[k] = value
			break
		}

		// not the last key, be deeper
		v, ok := cursor[k]
		if !ok {
			// create a new map
			deeper := make(map[string]any)
			cursor[k] = deeper
			cursor = deeper
			continue
		}

		// current exists, then check existing value type, if it's not map
		// that means duplicate keys, and at least one is not map instance.
		if cursor, ok = v.(map[string]any); !ok {
			log.Warnf("duplicate key: %v\n", strings.Join(keys[:i+1], "."))
			break
		}
	}
}

// genKey got the key of config.KeyValue pair.
// eg: namespace.ext with subKey got namespace.subKey
func genKey(ns, sub string) string {
	arr := strings.Split(ns, ".")
	if len(arr) == 1 {
		if ns == "" {
			return sub
		}
		return ns + "." + sub
	}
	if configFileformat(ns) != constant.DEFAULT {
		return strings.Join(arr[:len(arr)-1], ".") + "." + sub
	}
	return ns + "." + sub
}
