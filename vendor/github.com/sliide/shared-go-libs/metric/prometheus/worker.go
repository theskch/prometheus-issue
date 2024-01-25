package prometheus

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// WorkerCodeOK represents query OK that used for WorkerMetrics.Processed.
	WorkerCodeOK = "ok"

	// WorkerCodeFailed represents query failed that used for WorkerMetrics.Processed.
	WorkerCodeFailed = "failed"

	// WorkerCodeCanceled represents query canceled that used for WorkerMetrics.Processed.
	WorkerCodeCanceled = "canceled"
)

var (
	workerMessagesInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "worker_messages_in_flight",
			Help: "The current number of messages being processed by the worker",
		},
		[]string{"queue"},
	)

	workerMessagesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_messages_total",
			Help: "Total number of messages processed by the worker.",
		},
		[]string{"queue", "code"},
	)

	workerDeadLetterMessagesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_dead_letter_messages_total",
			Help: "Total number of messages sent to DLQ by the worker.",
		},
		[]string{"queue"},
	)

	workerMessagesProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "worker_messages_processing_duration_seconds",
			Help: "Duration of time from when the message is received from the queue until it is processed in seconds.",
		},
		[]string{"queue", "code"},
	)
)

// NewWorkerMetrics returns a metrics which used for monitoring worker processing.
func NewWorkerMetrics(queue string) *WorkerMetrics {
	return &WorkerMetrics{
		queue: queue,

		inFlight:     workerMessagesInFlight,
		total:        workerMessagesTotal,
		deadLetter:   workerDeadLetterMessagesTotal,
		durationSecs: workerMessagesProcessingDuration,

		t: time.Now(),
	}
}

// WorkerMetrics represents a collection of metrics to be registered on a Prometheus metrics registry.
type WorkerMetrics struct {
	queue string

	inFlight     *prometheus.GaugeVec
	total        *prometheus.CounterVec
	deadLetter   *prometheus.CounterVec
	durationSecs *prometheus.HistogramVec

	t     time.Time
	begin bool
}

// Begin the monitoring.
// Should be called when the worker reads a message from the queue.
func (m *WorkerMetrics) Begin() {
	if m.begin {
		return
	}
	m.t = time.Now()
	m.incInFlight()
	m.begin = true
}

// End the monitoring.
// Should be called when message processing has finished.
func (m *WorkerMetrics) End(code string) {
	if !m.begin {
		return
	}
	m.decInFlight()
	m.processed(code)
	m.observeProcessingTime(code, time.Since(m.t))
	m.begin = false
}

func (m WorkerMetrics) processed(code string) {
	m.total.With(m.labelWithCode(code)).Inc()
}

func (m WorkerMetrics) IncDeadLetter() {
	m.deadLetter.With(m.label()).Inc()
}

func (m WorkerMetrics) incInFlight() {
	m.inFlight.With(m.label()).Inc()
}

func (m WorkerMetrics) decInFlight() {
	m.inFlight.With(m.label()).Dec()
}

func (m WorkerMetrics) observeProcessingTime(code string, t time.Duration) {
	m.durationSecs.With(m.labelWithCode(code)).Observe(t.Seconds())
}

func (m WorkerMetrics) label() prometheus.Labels {
	return prometheus.Labels{
		"queue": strings.ToLower(m.queue),
	}
}

func (m WorkerMetrics) labelWithCode(code string) prometheus.Labels {
	label := m.label()
	label["code"] = strings.ToLower(code)

	return label
}
