package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	cacheSubsystem = "cache"
	success        = "success"
	failure        = "failure"
)

var (
	cacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: cacheSubsystem,
			Name:      "hits_total",
			Help:      "The total number of the cache hits",
		},
		[]string{"cache_name"},
	)
	cacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: cacheSubsystem,
			Name:      "misses_total",
			Help:      "The total number of the cache misses",
		},
		[]string{"cache_name"},
	)
	cacheSets = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: cacheSubsystem,
			Name:      "sets_total",
			Help:      "The total number of insertions to the cache",
		},
		[]string{"cache_name", "result"},
	)
)

// CacheMetricer describes the cache metrics.
type CacheMetricer interface {
	IncGet(err error)
	IncSet(err error)
}

// CacheMetrics is used to collect the cache metrics.
type CacheMetrics struct {
	cacheName string
}

// NewCacheMetrics returns the CacheMetrics instance.
func NewCacheMetrics(cacheName string) CacheMetrics {
	return CacheMetrics{
		cacheName: cacheName,
	}
}

// IncGet increments the number of the cache Get methods calls.
// Depending on the given error it increments the cacheHits or cacheMisses metric.
func (c CacheMetrics) IncGet(err error) {
	if err != nil {
		cacheMisses.WithLabelValues(c.cacheName).Inc()
	} else {
		cacheHits.WithLabelValues(c.cacheName).Inc()
	}
}

// IncSet increments the number of the cache Set method calls.
func (c CacheMetrics) IncSet(err error) {
	result := success
	if err != nil {
		result = failure
	}

	cacheSets.WithLabelValues(c.cacheName, result).Inc()
}
