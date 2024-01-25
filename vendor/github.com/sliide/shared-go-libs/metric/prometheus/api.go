package prometheus

import (
	"context"
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// apiCallOK represents the indicator that the API call was successful.
	apiCallOK = "ok"
	// apiCallFailed represents the indicator that the API call has failed.
	apiCallFailed = "failed"
	// apiCallFailed represents the indicator that the API call was canceled.
	apiCallCanceled = "canceled"
)

var (
	apiRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "third_party_api_requests_in_flight",
			Help: "The current number of third-party API requests being served.",
		},
		[]string{"api_name", "api_method"},
	)

	apiRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "third_party_api_requests_total",
			Help: "Total number of third-party API requests made and responded.",
		},
		[]string{"api_name", "api_method", "result"},
	)

	apiRequestsDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "third_party_api_request_duration_seconds",
			Help: "The request latencies of calling third-party API in seconds.",
		},
		[]string{"api_name", "api_method"},
	)
)

// APICallMetrics represents a collection of metrics to be registered on Prometheus metrics registry
// used to measure performances of the 3rd party APIs calls.
type APICallMetrics struct {
	apiName   string
	apiMethod string

	start time.Time
}

// NewAPICallMetrics retruns a metrics collection used to measure performances of the 3rd party API calls.
func NewAPICallMetrics(apiName, apiMethod string) APICallMetrics {
	return APICallMetrics{
		apiName:   apiName,
		apiMethod: apiMethod,
	}
}

// Begin starts monitoring the reqeust made to the API. The return function should be used
// when the response or error is received.
//
// Usage:
// endMonitoring := NewAPICallMetrics("someAPI", "someAPIMethod").Begin()
// result, err := someAPI.SomeAPIMethod()
// endMonitoring(err).
func (a APICallMetrics) Begin() func(err error) {
	a.start = time.Now()
	apiRequestsInFlight.With(a.labels()).Inc()

	return a.end
}

// end function is used to finalize API call monitoring.
func (a APICallMetrics) end(err error) {
	labels := a.labels()

	apiRequestsInFlight.With(labels).Dec()
	apiRequestsDuration.With(labels).Observe(time.Since(a.start).Seconds())
	apiRequestsTotal.With(a.labelsWithResult(err)).Inc()
}

// labels returns a map of predefined labels used for 3rd party API call monitoring.
func (a APICallMetrics) labels() prometheus.Labels {
	labels := prometheus.Labels{
		"api_name":   a.apiName,
		"api_method": a.apiMethod,
	}

	return labels
}

// labelsWithResult returns a map of predefined labels together with
// the result used for 3rd party API call monitoring.
func (a APICallMetrics) labelsWithResult(err error) prometheus.Labels {
	labels := a.labels()
	labels["result"] = apiCallResultFromError(err)

	return labels
}

// apiCallResultFromError returns the result of the API call from the provided error.
func apiCallResultFromError(err error) string {
	switch {
	case err == nil:
		return apiCallOK
	case errors.Is(err, context.Canceled):
		return apiCallCanceled
	default:
		return apiCallFailed
	}
}
