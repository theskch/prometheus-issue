package service

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ctxLoggerKey struct{}

func logPath(logger *logrus.Entry) func(next http.Handler) http.Handler {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logger.WithFields(logrus.Fields{
				"path":   r.URL.Path,
				"method": r.Method,
			})

			r = newRequestWithLogger(r, logger)
			next.ServeHTTP(w, r)
		})
	}
}

func logger(r *http.Request) *logrus.Entry {
	logger, ok := r.Context().Value(ctxLoggerKey{}).(*logrus.Entry)
	if !ok {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return logger
}

func newRequestWithLogger(r *http.Request, logger *logrus.Entry) *http.Request {
	ctx := context.WithValue(r.Context(), ctxLoggerKey{}, logger)
	r = r.WithContext(ctx)

	return r
}
