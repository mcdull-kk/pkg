package metric

import "github.com/prometheus/client_golang/prometheus"

type HistogramVecOpts struct {
	VectorOpts
	Buckets []float64
}

type HistogramVec interface {
	// Observe adds a single observation to the histogram.
	Observe(v int64, labels ...string)
}

// Histogram prom histogram collection.
type promHistogramVec struct {
	histogram *prometheus.HistogramVec
}

// NewHistogramVec new a histogram vec.
func NewHistogramVec(cfg *HistogramVecOpts) HistogramVec {
	if cfg == nil {
		return nil
	}
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
			Buckets:   cfg.Buckets,
		}, cfg.Labels)
	prometheus.MustRegister(vec)
	return &promHistogramVec{
		histogram: vec,
	}
}

// Timing adds a single observation to the histogram.
func (histogram *promHistogramVec) Observe(v int64, labels ...string) {
	histogram.histogram.WithLabelValues(labels...).Observe(float64(v))
}
