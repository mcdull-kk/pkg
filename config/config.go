package config

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
