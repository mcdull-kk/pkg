package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mcdull-kk/pkg/config"
)

var _ config.Source = (*file)(nil)

type file struct {
	path string
}

func NewSource(path string) config.Source {
	return &file{path: path}
}

func (f *file) Load() (kvs []*config.KeyValue, err error) {
	fi, err := os.Stat(f.path)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return f.loadDir(f.path)
	}
	kv, err := f.loadFile(f.path)
	if err != nil {
		return nil, err
	}
	kvs = append(kvs, kv)
	return
}

func (f *file) Watch() (config.Watcher, error) {
	return newWatcher(f)
}

func (f *file) loadDir(path string) (kvs []*config.KeyValue, err error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}
		kv, err := f.loadFile(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, kv)
	}
	return
}

func (f *file) loadFile(path string) (*config.KeyValue, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	format := ""
	if p := strings.Split(info.Name(), "."); len(p) > 1 {
		format = p[len(p)-1]
	}

	return &config.KeyValue{
		Key:    info.Name(),
		Format: format,
		Value:  data,
	}, nil
}
