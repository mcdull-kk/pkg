package file

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mcdull-kk/pkg/config"
	"github.com/mcdull-kk/pkg/log"
	"github.com/stretchr/testify/assert"
)

func Test_file(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test_config.json")
	defer os.Remove(path)
	log.Info(path)

	err := os.WriteFile(path, []byte(`{"name":"mucdull-kk","mucdull":{"name":"mucdull","age":2},"kk":{"name":"kk","age":1}}`), 0o666)
	assert.Nil(t, err)

	c := config.New(
		config.WithSource(
			NewSource(path),
		),
	)
	err = c.Load()
	assert.Nil(t, err)

	val := make(map[string]any)
	err = c.Scan(&val)
	assert.Nil(t, err)

	log.Info(val)
	log.Infow("mucdull.age", c.Value("mucdull.age"))

	f, err := os.OpenFile(path, os.O_RDWR, 0)
	assert.Nil(t, err)

	_, err = f.WriteString(`{"name":"mucdull-kk","mucdull":{"name":"mucdull","age":15},"kk":{"name":"kk","age":1}}`)
	assert.Nil(t, err)
	f.Close()

	time.Sleep(10 * time.Millisecond)

	log.Infow("mucdull.age", c.Value("mucdull.age"))
}
