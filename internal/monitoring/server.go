package monitoring

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sliide/shared-go-libs/metric/prometheus"
)

const (
	defaultShutdownTimeout = 5 * time.Second
)

type Server struct {
	server  *http.Server
	m       sync.Mutex
	serving int32
}

func NewServer(address string) *Server {
	r := chi.NewRouter()

	r.Get("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prometheus.Handler().ServeHTTP(w, r)
	}))

	s := &Server{
		server: &http.Server{
			Addr:    address,
			Handler: r,
		},
	}

	return s
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
