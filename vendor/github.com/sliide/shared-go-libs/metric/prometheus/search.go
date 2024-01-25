package prometheus

import (
	"sort"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// SearchMetrics represents a collection of metrics to be registered on a Prometheus metrics registry.
type SearchMetrics struct {
	providerName string
}

const (
	// SearchProviderCodeOK represents query OK that is used for ending monitoring of a Search Provider action.
	SearchProviderCodeOK = "ok"

	// SearchProviderCodeFailed represents query failed that is used for ending monitoring of a Search Provider action.
	SearchProviderCodeFailed = "failed"

	// SearchProviderCodeCanceled represents query canceled that is used for ending monitoring of a Search Provider action.
	SearchProviderCodeCanceled = "canceled"

	// SearchProviderCodeTimeout represents query timeout that is used for ending monitoring of a Search Provider action.
	SearchProviderCodeTimeout = "timeout"
)

var (
	searchIndexDocumentTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "search_index_document_total",
			Help: "Total number of index documents to Search Engine.",
		},
		[]string{"provider", "index", "code"},
	)

	searchIndexDocumentInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "search_index_document_in_flight",
			Help: "The current number of index documents to Search Engine.",
		},
		[]string{"provider", "index"},
	)

	searchIndexDocumentDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "search_index_document_duration_seconds",
			Help: "Duration of time it takes to index a new document to Search Engine.",
		},
		[]string{"provider", "index", "code"},
	)

	searchCreateIndexTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "search_create_index_total",
			Help: "Total number of index creations to Search Engine.",
		},
		[]string{"provider", "index", "code"},
	)

	searchCreateIndexInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "search_create_index_in_flight",
			Help: "The current number of index creations to Search Engine.",
		},
		[]string{"provider", "index"},
	)

	searchCreateIndexDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "search_create_index_duration_seconds",
			Help: "Duration of time it takes to create an index to Search Engine.",
		},
		[]string{"provider", "index", "code"},
	)

	searchGetDocumentTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "search_get_document_total",
			Help: "Total number of Search Engine get document by ID actions performed.",
		},
		[]string{"provider", "index", "code"},
	)

	searchGetDocumentInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "search_get_document_in_flight",
			Help: "The current number of get document by ID actions performed to Search Engine.",
		},
		[]string{"provider", "index"},
	)

	searchGetDocumentDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "search_get_document_duration_seconds",
			Help: "Duration of time it takes to get a document by ID from the Search Engine.",
		},
		[]string{"provider", "index", "code"},
	)

	searchQueryDocumentsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "search_query_documents_total",
			Help: "Total number of Search Engine document queries performed.",
		},
		[]string{"provider", "index", "code"},
	)

	searchQueryDocumentsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "search_query_documents_in_flight",
			Help: "The current number of document queries performed to Search Engine.",
		},
		[]string{"provider", "index"},
	)

	searchQueryDocumentsDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "search_query_documents_duration_seconds",
			Help: "Duration of time it takes to query documents from the Search Engine.",
		},
		[]string{"provider", "index", "code"},
	)

	pitRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "search_point_in_time_requests_total",
			Help: "Total number of requests to the Point In Time API",
		},
		[]string{"provider", "operation", "code"},
	)

	pitRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "search_point_in_time_requests_in_flight",
			Help: "The current number of requests to the Point In Time API",
		},
		[]string{"provider", "operation"},
	)

	pitRequestsDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "search_point_in_time_request_duration_seconds",
			Help: "Duration of requests made to Point In Time API",
		},
		[]string{"provider", "operation", "code"},
	)
)

// NewSearchMetrics returns a metrics which is used for monitoring Search provider interaction.
func NewSearchMetrics(providerName string) *SearchMetrics {
	return &SearchMetrics{providerName: providerName}
}

// BeginCreateIndex starts monitoring a CreateIndex action.
func (m *SearchMetrics) BeginCreateIndex(index string) func(code string) {
	start := time.Now()
	index = strings.ToLower(index)

	searchCreateIndexInFlight.WithLabelValues(m.providerName, index).Inc()

	return func(code string) {
		code = strings.ToLower(code)

		searchCreateIndexInFlight.WithLabelValues(m.providerName, index).Dec()
		searchCreateIndexTotal.WithLabelValues(m.providerName, index, code).Inc()
		searchCreateIndexDuration.WithLabelValues(m.providerName, index, code).Observe(time.Since(start).Seconds())
	}
}

// BeginIndexDocument starts monitoring an IndexDocument action.
func (m *SearchMetrics) BeginIndexDocument(index string) func(code string) {
	start := time.Now()
	index = strings.ToLower(index)

	searchIndexDocumentInFlight.WithLabelValues(m.providerName, index).Inc()

	return func(code string) {
		code = strings.ToLower(code)

		searchIndexDocumentInFlight.WithLabelValues(m.providerName, index).Dec()
		searchIndexDocumentTotal.WithLabelValues(m.providerName, index, code).Inc()
		searchIndexDocumentDuration.WithLabelValues(m.providerName, index, code).Observe(time.Since(start).Seconds())
	}
}

// BeginGetDocument starts monitoring a GetDocuments action.
func (m *SearchMetrics) BeginGetDocument(index string) func(code string) {
	start := time.Now()
	index = strings.ToLower(index)

	searchGetDocumentInFlight.WithLabelValues(m.providerName, index).Inc()

	return func(code string) {
		code = strings.ToLower(code)

		searchGetDocumentInFlight.WithLabelValues(m.providerName, index).Dec()
		searchGetDocumentTotal.WithLabelValues(m.providerName, index, code).Inc()
		searchGetDocumentDuration.WithLabelValues(m.providerName, index, code).Observe(time.Since(start).Seconds())
	}
}

// BeginSearchDocuments starts monitoring a QueryDocuments action.
func (m *SearchMetrics) BeginSearchDocuments(indexes []string) func(code string) {
	indexesStr := strFromSlice(indexes)
	start := time.Now()

	searchQueryDocumentsInFlight.WithLabelValues(m.providerName, indexesStr).Inc()

	return func(code string) {
		code = strings.ToLower(code)

		searchQueryDocumentsInFlight.WithLabelValues(m.providerName, indexesStr).Dec()
		searchQueryDocumentsTotal.WithLabelValues(m.providerName, indexesStr, code).Inc()
		searchQueryDocumentsDuration.WithLabelValues(m.providerName, indexesStr, code).Observe(time.Since(start).Seconds())
	}
}

// BeginPointInTimeRequest starts monitoring a PointInTime action.
func (m *SearchMetrics) BeginPointInTimeRequest(operation string) func(code string) {
	start := time.Now()
	operation = strings.ToLower(operation)

	pitRequestsInFlight.WithLabelValues(m.providerName, operation).Inc()

	return func(code string) {
		code = strings.ToLower(code)

		pitRequestsInFlight.WithLabelValues(m.providerName, operation).Dec()
		pitRequestsTotal.WithLabelValues(m.providerName, operation, code).Inc()
		pitRequestsDuration.WithLabelValues(m.providerName, operation, code).Observe(time.Since(start).Seconds())
	}
}

func strFromSlice(slice []string) string {
	result := make([]string, len(slice))
	for i, s := range slice {
		result[i] = strings.ToLower(s)
	}
	sort.Strings(result)

	return strings.Join(result, ",")
}
