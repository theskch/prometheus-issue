package prometheus

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// DynamoDBCodeOK represents query OK that is used for ending monitoring of an DynamoDB action.
	DynamoDBCodeOK = "ok"

	// DynamoDBCodeFailed represents query failed that is used for ending monitoring of an DynamoDB action.
	DynamoDBCodeFailed = "failed"

	// DynamoDBCodeCanceled represents query canceled that is used for ending monitoring of an DynamoDB action.
	DynamoDBCodeCanceled = "canceled"
)

var (
	dynamoPut = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dynamo_put_total",
			Help: "Total number of objects put to DynamoDB.",
		},
		[]string{"table", "code"},
	)

	dynamoPutDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "dynamo_put_duration_seconds",
			Help: "Duration of time it takes to put a new object to DynamoDB.",
		},
		[]string{"table", "code"},
	)

	dynamoGet = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dynamo_get_total",
			Help: "Total number of objects got from DynamoDB.",
		},
		[]string{"table", "code"},
	)

	dynamoGetDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "dynamo_get_duration_seconds",
			Help: "Duration of time it takes to get an object from DynamoDB.",
		},
		[]string{"table", "code"},
	)
)

// NewDynamoDBMetrics returns a metrics which is used for monitoring DynamoDB interaction.
func NewDynamoDBMetrics(table string) *DynamoDBMetrics {
	return &DynamoDBMetrics{
		table: table,
	}
}

// DynamoDBMetrics represents a collection of metrics to be registered on a Prometheus metrics registry.
type DynamoDBMetrics struct {
	table string
}

type dynamoMetrics struct {
	total        *prometheus.CounterVec
	durationSecs *prometheus.HistogramVec
	table        string

	start time.Time
}

// BeginPut starts monitoring a DynamoDBMetrics put action.
func (m *DynamoDBMetrics) BeginPut() func(code string) {
	metric := &dynamoMetrics{
		table:        m.table,
		total:        dynamoPut,
		durationSecs: dynamoPutDuration,
	}
	metric.begin()

	return metric.end
}

// BeginGet starts monitoring a DynamoDBMetrics get action.
func (m *DynamoDBMetrics) BeginGet() func(code string) {
	metric := &dynamoMetrics{
		table:        m.table,
		total:        dynamoGet,
		durationSecs: dynamoGetDuration,
	}

	metric.begin()

	return metric.end
}

func (m *dynamoMetrics) begin() {
	m.start = time.Now()
}

func (m *dynamoMetrics) end(code string) {
	m.total.With(m.labelWithCode(code)).Inc()
	m.durationSecs.With(m.labelWithCode(code)).Observe(time.Since(m.start).Seconds())
}

func (m *dynamoMetrics) labelWithCode(code string) prometheus.Labels {
	label := m.label()
	label["code"] = strings.ToLower(code)

	return label
}

func (m *dynamoMetrics) label() prometheus.Labels {
	return prometheus.Labels{
		"table": strings.ToLower(m.table),
	}
}
