package redis

import (
	"errors"
	"io"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestRedis_Exists(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		cli, err := client.Client()
		assert.NotNil(t, err)
		_, err = New(client.Addr, badType()).Client()
		assert.NotNil(t, err)

		v, err := cli.Exists("a").Result()
		assert.Nil(t, err)
		assert.Zero(t, v)
		_, err = cli.Set("a", "b", 0).Result()
		assert.Nil(t, err)
		v, err = cli.Exists("a").Result()
		assert.Nil(t, err)
		assert.NotZero(t, v)
	})
}

func badType() Option {
	return func(r *Redis) {
		r.Type = "bad"
	}
}

func runOnRedis(t *testing.T, fn func(client *Redis)) {
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer func() {
		client, err := _clientManager.GetResource(s.Addr(), func() (io.Closer, error) {
			return nil, errors.New("should already exist")
		})
		if err != nil {
			t.Error(err)
		}

		if client != nil {
			client.Close()
		}
	}()
	fn(NewRedis(s.Addr(), NodeType))
}
