package monitoring

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	apiRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "api_requests_in_flight",
			Help: "Number of api requests in flight",
		},
		[]string{"ID"})

	apiRequestDuration = promauto.NewHistogramVec(

		prometheus.HistogramOpts{
			Name: "api_request_duration_seconds",
			Help: "Duration of api requests in seconds",
		},
		[]string{"ID"},
	)
)

type APIRequestMetric struct {
	id    string
	start time.Time
}

func NewAPIRequestMetric(id int) APIRequestMetric {
	return APIRequestMetric{
		id:    fmt.Sprintf("%d", id),
		start: time.Now(),
	}
}

func (m APIRequestMetric) Begin() func() {
	apiRequestsInFlight.WithLabelValues(m.id).Inc()

	return func() {
		apiRequestsInFlight.WithLabelValues(m.id).Dec()
		apiRequestDuration.WithLabelValues(m.id).Observe(time.Since(m.start).Seconds())
	}
}
