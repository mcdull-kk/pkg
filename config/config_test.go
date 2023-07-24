package config

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mcdull-kk/pkg/codec"
	"github.com/stretchr/testify/assert"
)

func Test_defaultDecoder(t *testing.T) {
	tests := []struct {
		keyValue *KeyValue
		want     map[string]any
	}{
		{
			keyValue: &KeyValue{
				Key:    "service",
				Value:  []byte("config"),
				Format: "",
			},
			want: map[string]any{"service": []byte("config")},
		},
		{
			keyValue: &KeyValue{
				Key:    "service.name.alias",
				Value:  []byte("2233333"),
				Format: "",
			},
			want: map[string]interface{}{
				"service": map[string]interface{}{
					"name": map[string]interface{}{
						"alias": []byte("2233333"),
					},
				},
			},
		},
		{
			keyValue: &KeyValue{
				Key:    "service.name.alias",
				Value:  []byte(`{"name":"alias"}`),
				Format: "json",
			},
			want: map[string]interface{}{
				"name": "alias",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.keyValue.Key, func(t *testing.T) {
			target := make(map[string]any)
			defaultDecoder(tt.keyValue, target)
			assert.Equal(t, tt.want, target)
		})
	}
}

func Test_defaultResolver(t *testing.T) {
	var ()

	data := map[string]any{
		"foo": map[string]any{
			"bar": map[string]any{
				"notexist": "${NOTEXIST:100}",
				"port":     "${PORT:8081}",
				"count":    "${COUNT:0}",
				"enable":   "${ENABLE:false}",
				"rate":     "${RATE}",
				"empty":    "${EMPTY:foobar}",
				"url":      "${URL:http://example.com}",
				"array": []any{
					"${PORT}",
					map[string]any{"foobar": "${NOTEXIST:8081}"},
				},
				"value1": "${test.value}",
				"value2": "$PORT",
				"value3": "abc${PORT}foo${COUNT}bar",
				"value4": "${foo${bar}}",
			},
		},
		"test": map[string]any{
			"value": "foobar",
		},
		"PORT":   "8080",
		"COUNT":  "10",
		"ENABLE": "true",
		"RATE":   "0.9",
		"EMPTY":  "",
	}

	tests := []struct {
		path string
		want any
	}{
		{
			path: "foo.bar.notexist",
			want: 100,
		},
		{
			path: "foo.bar.port",
			want: "8080",
		},
		{
			path: "foo.bar.count",
			want: 10,
		},
		{
			path: "foo.bar.enable",
			want: true,
		},
		{
			path: "foo.bar.rate",
			want: 0.9,
		},
		{
			path: "foo.bar.empty",
			want: "",
		},
		{
			path: "foo.bar.url",
			want: "http://example.com",
		},
		{
			path: "foo.bar.array",
			want: []any{"8080", map[string]any{"foobar": "8081"}},
		},
		{
			path: "foo.bar.value1",
			want: "foobar",
		},
		{
			path: "foo.bar.value2",
			want: "$PORT",
		},
		{
			path: "foo.bar.value3",
			want: "abc8080foo10bar",
		},
		{
			path: "foo.bar.value4",
			want: "}",
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			err := defaultResolver(data)
			if err != nil {
				t.Fatal(err)
			}
			rd := reader{
				values: data,
			}
			if vv, ok := rd.Value(test.path); ok {
				v := vv.Load()
				switch test.want.(type) {
				case int:
					if !reflect.DeepEqual(test.want.(int), int(codec.Int(v))) {
						t.Fatal("want is not equal to actual")
					}
				case string:
					if !reflect.DeepEqual(test.want, codec.Repr(v)) {
						t.Fatal("want is not equal to actual")
					}
				case bool:
					if !reflect.DeepEqual(test.want, codec.Bool(v)) {
						t.Fatal("want is not equal to actual")
					}
				case float64:
					if !reflect.DeepEqual(test.want, codec.Float(v)) {
						t.Fatal("want is not equal to actual")
					}
				default:
					if !reflect.DeepEqual(test.want, v) {
						t.Logf("want: %#v, actural: %#v", test.want, v)
						t.Fail()
					}
				}
				if err != nil {
					t.Error(err)
				}
			} else {
				t.Error("value path not found")
			}
		})
	}
}

func TestExpand(t *testing.T) {
	tests := []struct {
		input   string
		mapping func(string) string
		want    string
	}{
		{
			input: "${a}",
			mapping: func(s string) string {
				return strings.ToUpper(s)
			},
			want: "A",
		},
		{
			input: "a",
			mapping: func(s string) string {
				return strings.ToUpper(s)
			},
			want: "a",
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, expand(tt.input, tt.mapping))
	}
}
