package metric

import (
	"errors"
	"fmt"
)

type VectorOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

const (
	_businessSubsystemCount     = "count"
	_businessSubSystemGauge     = "gauge"
	_businessSubSystemHistogram = "histogram"
)

var (
	_defaultBuckets = []float64{5, 10, 25, 50, 100, 250, 500}
)

// NewMetricCount business Metric count vec.
// name or labels should not be empty.
func NewMetricCount(name string, labels ...string) CounterVec {
	if name == "" || len(labels) == 0 {
		panic(errors.New("stat:metric business count metric name should not be empty or labels length should be greater than zero"))
	}
	return NewCounterVec(&CounterVecOpts{
		Subsystem: _businessSubsystemCount,
		Name:      name,
		Labels:    labels,
		Help:      fmt.Sprintf("metric count %s", name),
	})
}

// NewMetricGauge business Metric gauge vec.
// name or labels should not be empty.
func NewMetricGauge(name string, labels ...string) GaugeVec {
	if name == "" || len(labels) == 0 {
		panic(errors.New("stat:metric business gauge metric name should not be empty or labels length should be greater than zero"))
	}
	return NewGaugeVec(&GaugeVecOpts{
		Subsystem: _businessSubSystemGauge,
		Name:      name,
		Labels:    labels,
		Help:      fmt.Sprintf("metric gauge %s", name),
	})
}

// NewBusinessMetricHistogram business Metric histogram vec.
// name or labels should not be empty.
func NewBusinessMetricHistogram(name string, buckets []float64, labels ...string) HistogramVec {
	if name == "" || len(labels) == 0 {
		panic(errors.New("stat:metric business histogram metric name should not be empty or labels length should be greater than zero"))
	}
	if len(buckets) == 0 {
		buckets = _defaultBuckets
	}
	return NewHistogramVec(&HistogramVecOpts{
		VectorOpts: VectorOpts{
			Subsystem: _businessSubSystemHistogram,
			Name:      name,
			Labels:    labels,
			Help:      fmt.Sprintf("metric histogram %s", name),
		},
		Buckets: buckets,
	})
}
