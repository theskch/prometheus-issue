package service

import (
	"net/http"

	"github.com/theskch/prometheus-issue/internal/monitoring"
)

type app struct{}

func (a app) Info(w http.ResponseWriter, r *http.Request, id int) {
	metric := monitoring.NewAPIRequestMetric(id)

	end := metric.Begin()
	defer end()

	log := logger(r)
	log.Info("Request received")
}

func (a app) Ping(w http.ResponseWriter, r *http.Request) {
	_ = renderRawJSON(w, http.StatusOK, nil)
}
