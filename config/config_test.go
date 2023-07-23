package config

import (
	"testing"

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

}
