package service

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/theskch/prometheus-issue/pkg/api"
)

const (
	defaultReadTimeout     = 10 * time.Second
	defaultShutdownTimeout = 5 * time.Second
)

var (
	notFoundJSONResponseBody = func() []byte {
		b, _ := json.Marshal(&api.Error{
			Error: "not-found",
		})

		return b
	}()

	methodNotAllowedResponseBody = func() []byte {
		b, _ := json.Marshal(&api.Error{
			Error: "method-not-allowed",
		})

		return b
	}()
)

type Server struct {
	server  *http.Server
	serving int32
	m       sync.Mutex
}

func NewServer(address string) *Server {
	r := chi.NewRouter()

	r.NotFound(notFoundHandler)
	r.MethodNotAllowed(methodNotAllowedHandler)

	logger := logrus.NewEntry(logrus.StandardLogger())
	r.Use(logPath(logger))

	serverOptions := api.ChiServerOptions{
		BaseURL:          "/v1",
		BaseRouter:       r,
		ErrorHandlerFunc: errorHandler,
	}

	api.HandlerWithOptions(app{}, serverOptions)

	return &Server{
		server: &http.Server{
			Addr:        address,
			Handler:     r,
			ReadTimeout: defaultReadTimeout,
		},
	}
}

func (s *Server) ListenAndServe() error {
	s.m.Lock()
	defer s.m.Unlock()

	atomic.StoreInt32(&s.serving, 1)
	defer func() {
		atomic.StoreInt32(&s.serving, 0)
	}()

	return s.server.ListenAndServe()
}

func (s *Server) GracefulStop() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func (s *Server) Serving() bool {
	return atomic.LoadInt32(&s.serving) == 1
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	_ = renderRawJSON(w, http.StatusNotFound, notFoundJSONResponseBody)
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	_ = renderRawJSON(w, http.StatusMethodNotAllowed, methodNotAllowedResponseBody)
}

func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse, _ := json.Marshal(&api.Error{
		Error: err.Error(),
	})

	_ = renderRawJSON(w, http.StatusBadRequest, errorResponse)
}

func renderRawJSON(w http.ResponseWriter, statusCode int, payload []byte) (err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload != nil {
		_, err = w.Write(payload)
	}

	return
}
