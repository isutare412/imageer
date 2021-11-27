package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/isutare412/imageer/api/pkg/config"
)

type server struct {
	server *http.Server
}

func (s *server) Start() <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("on server.Start: %v", err)
			return
		}
		log.Info("HTTP server finished serving")
	}()

	log.Infof("HTTP server started on %s", s.server.Addr)
	return errChan
}

func (s *server) Shutdown() {
	log.Info("HTTP server shutdown start")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown HTTP server: %v", err)
		return
	}
	log.Info("HTTP server shutdown finished successfully")
}

func New(cfg *config.HttpConfig) *server {
	r := mux.NewRouter()

	// TODO: Attach middlewares

	apiV1 := r.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/greeting/{name}", getGreeting)

	return &server{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Handler: r,
		},
	}
}
