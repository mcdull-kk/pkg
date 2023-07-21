package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mcdull-kk/pkg/codec"
)

type (
	KeyValue struct {
		Key    string
		Value  []byte
		Format string
	}

	Source interface {
		Load() ([]*KeyValue, error)
		Watch() (Watcher, error)
	}

	Watcher interface {
		Next() ([]*KeyValue, error)
		Stop() error
	}
)

type (
	Decoder  func(*KeyValue, map[string]any) error
	Resolver func(map[string]any) error
	Option   func(*options)

	options struct {
		sources  []Source
		decoder  Decoder
		resolver Resolver
	}
)

func WithSource(s ...Source) Option {
	return func(o *options) {
		o.sources = s
	}
}

func WithDecoder(d Decoder) Option {
	return func(o *options) {
		o.decoder = d
	}
}

func WithResolver(r Resolver) Option {
	return func(o *options) {
		o.resolver = r
	}
}

// defaultDecoder decode config from source KeyValue
// to target map[string]any using src.Format codec.
func defaultDecoder(src *KeyValue, target map[string]any) error {
	if src.Format == "" {
		// expand key "aaa.bbb" into map[aaa]map[bbb]any
		keys := strings.Split(src.Key, ".")
		for i, k := range keys {
			if i == len(keys)-1 {
				target[k] = src.Value
			} else {
				sub := make(map[string]any)
				target[k] = sub
				target = sub
			}
		}
		return nil
	}
	if code := codec.GetCodec(src.Format); code != nil {
		return code.Unmarshal(src.Value, &target)
	}
	return fmt.Errorf("unsupported key: %s format: %s", src.Key, src.Format)
}

// defaultResolver resolve placeholder in map value,
// placeholder format in ${key:default}.
func defaultResolver(input map[string]any) error {
	mapper := func(name string) string {
		args := strings.SplitN(strings.TrimSpace(name), ":", 2) //nolint:gomnd
		if v, has := readValue(input, args[0]); has {
			return codec.Repr(v.Load())
		} else if len(args) > 1 { // default value
			return args[1]
		}
		return ""
	}

	var resolve func(map[string]any) error
	resolve = func(sub map[string]any) error {
		for k, v := range sub {
			switch vt := v.(type) {
			case string:
				sub[k] = expand(vt, mapper)
			case map[string]any:
				if err := resolve(vt); err != nil {
					return err
				}
			case []any:
				for i, iface := range vt {
					switch it := iface.(type) {
					case string:
						vt[i] = expand(it, mapper)
					case map[string]any:
						if err := resolve(it); err != nil {
							return err
						}
					}
				}
				sub[k] = vt
			}
		}
		return nil
	}
	return resolve(input)
}

func expand(s string, mapping func(string) string) string {
	r := regexp.MustCompile(`\${(.*?)}`)
	re := r.FindAllStringSubmatch(s, -1)
	for _, i := range re {
		if len(i) == 2 { //nolint:gomnd
			s = strings.ReplaceAll(s, i[0], mapping(i[1]))
		}
	}
	return s
}
