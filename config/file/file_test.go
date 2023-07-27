package file

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mcdull-kk/pkg/codec"
	"github.com/mcdull-kk/pkg/config"
	"github.com/mcdull-kk/pkg/log"
	"github.com/stretchr/testify/assert"
)

func Test_file(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test_config.json")
	defer os.Remove(path)
	log.Info(path)

	data := []byte(`{"name":"mucdull-kk","mucdull":{"name":"mucdull","age":2},"kk":{"name":"kk","age":1}}`)
	err := os.WriteFile(path, data, 0o666)
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

	assert.Equal(t, 2, int(codec.Int(c.Value("mucdull.age").Load())))

	updateData := `{"name":"mucdull-kk","mucdull":{"name":"mucdull","age":15},"kk":{"name":"kk","age":1}}`
	err = updateFileData(path, updateData)
	assert.Nil(t, err)

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 15, int(codec.Int(c.Value("mucdull.age").Load())))
}

func updateFileData(file string, updateData string) (err error) {
	f, err := os.OpenFile(file, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(updateData)
	return
}

func Test_source(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test_config")
	file := filepath.Join(path, "test.json")
	data := []byte(`{"name":"mucdull-kk","mucdull":{"name":"mucdull","age":2},"kk":{"name":"kk","age":1}}`)
	defer os.Remove(path)
	err := os.Mkdir(path, 0o700)
	assert.Nil(t, err)
	err = os.WriteFile(file, data, 0o666)
	assert.Nil(t, err)

	testPaths := []string{path, file}
	for _, tp := range testPaths {
		log.Info(tp)
		s := NewSource(tp)
		kvs, err := s.Load()
		assert.Nil(t, err)
		assert.Equal(t, string(data), string(kvs[0].Value))
	}
}

func Test_watch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test_config")
	file := filepath.Join(path, "test.json")
	data := []byte(`{"name":"mucdull-kk","mucdull":{"name":"mucdull","age":2},"kk":{"name":"kk","age":1}}`)
	defer os.Remove(path)
	err := os.Mkdir(path, 0o700)
	assert.Nil(t, err)
	err = os.WriteFile(file, data, 0o666)
	assert.Nil(t, err)

	testPaths := []struct {
		p    string
		want string
	}{
		{
			p:    path,
			want: `{"name":"mucdull-kk","mucdull":{"name":"mucdull","age":15},"kk":{"name":"kk","age":1}}`,
		},
		{
			p:    file,
			want: `{"name":"mucdull-kk","mucdull":{"name":"mucdull","age":20},"kk":{"name":"kk","age":1}}`,
		},
	}

	for _, tp := range testPaths {
		p := tp.p
		log.Info(p)
		f := p

		s := NewSource(p)
		watch, err := s.Watch()
		assert.Nil(t, err)

		fi, err := os.Stat(p)
		assert.Nil(t, err)
		if fi.IsDir() {
			fs, err := os.ReadDir(p)
			assert.Nil(t, err)
			f = filepath.Join(p, fs[0].Name())
		}

		err = updateFileData(f, tp.want)
		assert.Nil(t, err)

		kvs, err := watch.Next()
		assert.Nil(t, err)
		assert.Equal(t, tp.want, string(kvs[0].Value))
	}
}
