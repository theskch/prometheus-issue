package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/theskch/prometheus-issue/internal/monitoring"
	"github.com/theskch/prometheus-issue/internal/service"
)

const (
	monitoringServerAddress = ":9090"
	serviceServerAddress    = ":8080"
)

func main() {
	monitoringServer := monitoring.NewServer(monitoringServerAddress)
	log := logrus.NewEntry(logrus.StandardLogger())

	log.Info("Starting monitoring server")
	go func() {
		if err := monitoringServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Error("Error while starting monitoring server")
		}
	}()

	log.Info("Starting service server")
	serviceServer := service.NewServer(serviceServerAddress)
	go func() {
		if err := serviceServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Error("Error while starting service server")
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	logrus.Info("Stopping service server")
	if err := serviceServer.GracefulStop(); err != nil {
		logrus.WithError(err).Fatal("Failed to gracefully stop the http server")
	}

	logrus.Info("Stopping monitoring server")
	if err := monitoringServer.GracefulStop(); err != nil {
		logrus.WithError(err).Warning("Failed to gracefully stop the monitoring server")
	}
}
