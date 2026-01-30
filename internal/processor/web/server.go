package web

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/isutare412/imageer/internal/processor/metric"
)

type Server struct {
	cfg    Config
	server *http.Server
}

func NewServer(cfg Config) *Server {
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.HandlerFor(metric.Gatherer(), promhttp.HandlerOpts{}))

	return &Server{
		cfg: cfg,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: r,
		},
	}
}

func (s *Server) Run() <-chan error {
	errs := make(chan error, 1)
	go func() {
		defer close(errs)

		slog.Info("Starting web server", "port", s.cfg.Port)
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errs <- fmt.Errorf("failed to start server: %w", err)
			return
		}
	}()

	return errs
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.Background())
}
