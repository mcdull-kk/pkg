package redis

import (
	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
)

type RedisPoolStats interface {
	PoolStats() *redis.PoolStats
}

type collector struct {
	ps RedisPoolStats

	hits     *prometheus.Desc // number of times free connection was found in the pool
	misses   *prometheus.Desc // number of times free connection was NOT found in the pool
	timeouts *prometheus.Desc // number of times a wait timeout occurred

	totalConns *prometheus.Desc // number of total connections in the pool
	idleConns  *prometheus.Desc // number of idle connections in the pool
	staleConns *prometheus.Desc // number of stale connections removed from the pool
}

func NewCollector(ps RedisPoolStats, labels prometheus.Labels) prometheus.Collector {
	fqName := func(name string) string {
		return "go_redis_" + name
	}
	return &collector{
		ps: ps,
		hits: prometheus.NewDesc(
			fqName("total_hits"),
			"number of times free connection was found in the pool",
			nil,
			labels,
		),
		misses: prometheus.NewDesc(
			fqName("total_misses"),
			"number of times free connection was NOT found in the pool",
			nil,
			labels,
		),
		timeouts: prometheus.NewDesc(
			fqName("total_timeouts"),
			"number of times a wait timeout occurred",
			nil,
			labels,
		),
		totalConns: prometheus.NewDesc(
			fqName("total_connections"),
			"number of total connections in the pool",
			nil,
			labels,
		),
		idleConns: prometheus.NewDesc(
			fqName("idle_connections"),
			"number of idle connections in the pool",
			nil,
			labels,
		),
		staleConns: prometheus.NewDesc(
			fqName("stale_connections"),
			"number of stale connections removed from the pool",
			nil,
			labels,
		),
	}

}

// Describe implements Collector.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.hits
	ch <- c.misses
	ch <- c.timeouts
	ch <- c.totalConns
	ch <- c.idleConns
	ch <- c.staleConns
}

// Collect implements Collector.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	stats := c.ps.PoolStats()
	ch <- prometheus.MustNewConstMetric(c.hits, prometheus.GaugeValue, float64(stats.Hits))
	ch <- prometheus.MustNewConstMetric(c.misses, prometheus.GaugeValue, float64(stats.Misses))
	ch <- prometheus.MustNewConstMetric(c.timeouts, prometheus.GaugeValue, float64(stats.Timeouts))
	ch <- prometheus.MustNewConstMetric(c.totalConns, prometheus.GaugeValue, float64(stats.TotalConns))
	ch <- prometheus.MustNewConstMetric(c.idleConns, prometheus.GaugeValue, float64(stats.IdleConns))
	ch <- prometheus.MustNewConstMetric(c.staleConns, prometheus.GaugeValue, float64(stats.StaleConns))
}
