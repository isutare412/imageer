package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/isutare412/imageer/api-server/api"
	"github.com/isutare412/imageer/api-server/pkg/config"
	"github.com/isutare412/imageer/api-server/pkg/core/job"
	"github.com/isutare412/imageer/api-server/pkg/core/user"
)

type server struct {
	server *http.Server
	done   chan struct{}
}

func (s *server) Start(ctx context.Context) <-chan error {
	errNotify := make(chan error)
	go func() {
		defer close(errNotify)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errNotify <- fmt.Errorf("on http listen: %w", err)
		}
	}()
	log.Infof("HTTP server started on %s", s.server.Addr)

	go func() {
		<-ctx.Done()
		s.shutdown()
	}()

	return errNotify
}

func (s *server) shutdown() {
	log.Info("HTTP server shutdown start")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	defer close(s.done)

	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorf("Failed to shutdown HTTP server: %v", err)
	}

	log.Info("HTTP server shutdown finished successfully")
}

func (s *server) Done() <-chan struct{} {
	return s.done
}

func NewServer(cfg *config.HttpConfig, jSvc job.Service, uSvc user.Service) *server {
	r := mux.NewRouter()

	r.Use(logRequest, allowCORS)

	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler).Methods("GET")

	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	apiV1.HandleFunc("/greetings/{name}", getGreeting(jSvc)).Methods("GET")

	apiV1.HandleFunc("/users", createUser(uSvc)).Methods("POST")
	apiV1.HandleFunc("/users/{id}", getUser(uSvc)).Methods("GET")

	return &server{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Handler: r,
		},
		done: make(chan struct{}),
	}
}
