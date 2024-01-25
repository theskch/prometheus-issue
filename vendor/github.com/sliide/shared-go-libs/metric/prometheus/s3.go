package prometheus

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// S3CodeOK represents query OK that is used for ending monitoring of an S3 action.
	S3CodeOK = "ok"

	// S3CodeFailed represents query failed that is used for ending monitoring of an S3 action.
	S3CodeFailed = "failed"

	// S3CodeCanceled represents query canceled that is used for ending monitoring of an S3 action.
	S3CodeCanceled = "canceled"
)

var (
	s3ObjectsPut = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "s3_objects_put_total",
			Help: "Total number of objects put to S3",
		},
		[]string{"bucket", "code"},
	)

	s3ObjectsPutDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "s3_objects_put_duration_seconds",
			Help: "Duration of time it takes to put a new object to S3.",
		},
		[]string{"bucket", "code"},
	)

	s3ObjectsGet = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "s3_objects_get_total",
			Help: "Total number of objects got from S3.",
		},
		[]string{"bucket", "code"},
	)

	s3ObjectsGetDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "s3_objects_got_duration_seconds",
			Help: "Duration of time it takes to get an object from S3.",
		},
		[]string{"bucket", "code"},
	)
)

// NewS3Metrics returns a metrics which used for monitoring S3 interaction.
func NewS3Metrics(bucket string) *S3Metrics {
	return &S3Metrics{
		bucket: bucket,
	}
}

// S3Metrics represents a collection of metrics to be registered on a Prometheus metrics registry.
type S3Metrics struct {
	bucket string
}

type actionMetrics struct {
	total        *prometheus.CounterVec
	durationSecs *prometheus.HistogramVec
	bucket       string

	start time.Time
}

// BeginPut starts monitoring a S3 put action.
func (m *S3Metrics) BeginPut() func(code string) {
	metric := &actionMetrics{
		total:        s3ObjectsPut,
		durationSecs: s3ObjectsPutDuration,
	}
	metric.begin()

	return metric.end
}

// BeginGet starts monitoring a S3 get action.
func (m *S3Metrics) BeginGet() func(code string) {
	metric := &actionMetrics{
		total:        s3ObjectsGet,
		durationSecs: s3ObjectsGetDuration,
	}

	metric.begin()

	return metric.end
}

func (m *actionMetrics) begin() {
	m.start = time.Now()
}

func (m *actionMetrics) end(code string) {
	m.total.With(m.labelWithCode(code)).Inc()
	m.durationSecs.With(m.labelWithCode(code)).Observe(time.Since(m.start).Seconds())
}

func (m *actionMetrics) labelWithCode(code string) prometheus.Labels {
	label := m.label()
	label["code"] = strings.ToLower(code)

	return label
}

func (m *actionMetrics) label() prometheus.Labels {
	return prometheus.Labels{
		"bucket": strings.ToLower(m.bucket),
	}
}
