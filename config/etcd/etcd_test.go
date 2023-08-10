package etcd

import (
	"context"
	"testing"
	"time"

	"github.com/mcdull-kk/pkg/codec"
	"github.com/mcdull-kk/pkg/config"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func Test_etcd(t *testing.T) {
	path := "/mcdull-kk/test/config"
	source, client := NewSource(
		WithEndpoints([]string{"127.0.0.1:2379"}),
		WithDialTimeout(time.Second),
		WithDialOptions([]grpc.DialOption{grpc.WithBlock()}),
		WithPath(path),
	)

	// 保证先有值
	// _, err := client.Put(context.Background(), path, "test config")
	// assert.Nil(t, err)

	c := config.New(
		config.WithSource(
			source,
		),
	)
	defer c.Close()

	tests := []struct {
		want string
	}{
		{
			want: "old config",
		},
		{
			want: "new config",
		},
	}

	err := c.Load()
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			_, err = client.Put(context.Background(), path, tt.want)
			time.Sleep(20 * time.Millisecond)
			v := c.Value(path)
			assert.Nil(t, v)
			assert.Equal(t, tt.want, codec.Repr(v.Load()))
		})
	}

	_, err = client.Delete(context.Background(), path)
	assert.Nil(t, err)
}
