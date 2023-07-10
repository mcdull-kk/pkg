package redis

import (
	"crypto/tls"
	"io"

	red "github.com/go-redis/redis"
	manager "github.com/mcdull-kk/pkg/stores"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultDatabase = 0
	maxRetries      = 3
	idleConns       = 8
)

var (
	_clientManager  = manager.NewResourceManager()
	_clusterManager = manager.NewResourceManager()
)

func getClient(r *Redis) (*red.Client, error) {
	val, err := _clientManager.GetResource(r.Addr, func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewClient(&red.Options{
			Addr:         r.Addr,
			Password:     r.Pass,
			DB:           defaultDatabase,
			MaxRetries:   maxRetries,
			MinIdleConns: idleConns,
			TLSConfig:    tlsConfig,
		})
		store.WrapProcess(process)
		if r.collector {
			prometheus.MustRegister(NewCollector(store, prometheus.Labels{
				"addr": r.Addr,
			}))
		}
		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.Client), nil
}

func getCluster(r *Redis) (*red.ClusterClient, error) {
	val, err := _clusterManager.GetResource(r.Addr, func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewClusterClient(&red.ClusterOptions{
			Addrs:        []string{r.Addr},
			Password:     r.Pass,
			MaxRetries:   maxRetries,
			MinIdleConns: idleConns,
			TLSConfig:    tlsConfig,
		})
		store.WrapProcess(process)
		if r.collector {
			prometheus.MustRegister(NewCollector(store, prometheus.Labels{
				"addr": r.Addr,
			}))
		}

		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.ClusterClient), nil
}
