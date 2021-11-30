package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/isutare412/imageer/api/internal/config"
)

type server struct {
	server *http.Server
}

func (s *server) Start() <-chan error {
	errors := make(chan error)

	go func() {
		defer close(errors)

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errors <- fmt.Errorf("on server.Start: %v", err)
			return
		}
		log.Info("HTTP server finished serving")
	}()

	log.Infof("HTTP server started on %s", s.server.Addr)
	return errors
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

	r.Use(logRequest, allowCORS)

	apiV1 := r.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/greeting/{name}", getGreeting).Methods("GET")

	return &server{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Handler: r,
		},
	}
}
