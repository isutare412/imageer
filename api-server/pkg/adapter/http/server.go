package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/isutare412/imageer/api-server/api"
	"github.com/isutare412/imageer/api-server/pkg/config"
	"github.com/isutare412/imageer/api-server/pkg/core/auth"
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
			errNotify <- err
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
		log.Errorf("failed to shutdown HTTP server: %v", err)
	}

	log.Info("HTTP server shutdown finished successfully")
}

func (s *server) Done() <-chan struct{} {
	return s.done
}

func NewServer(
	cfg *config.HttpConfig, jSvc job.Service, uSvc user.Service, authSvc auth.Service,
) *server {
	injectSession := injectSession(authSvc)
	checkAdmin := checkAdmin(authSvc)

	r := mux.NewRouter()
	r.Use(logRequest)

	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler).Methods("GET")
	r.HandleFunc("/signIn", signIn(uSvc, authSvc)).Methods("POST")
	r.HandleFunc("/signOut", signOut(uSvc, authSvc)).Methods("GET")

	apiBase := r.PathPrefix("/api/v1").Subrouter()

	apiAuth := apiBase.NewRoute().Subrouter()
	apiAuth.Use(injectSession)

	apiAdmin := apiBase.NewRoute().Subrouter()
	apiAdmin.Use(injectSession, checkAdmin)

	apiAuth.HandleFunc("/greetings/{name}", getGreeting(jSvc)).Methods("GET")

	apiAuth.HandleFunc("/users", getUser(uSvc, authSvc)).Methods("GET")
	apiBase.HandleFunc("/users", createUser(uSvc)).Methods("POST")
	apiAdmin.HandleFunc("/users/{id}", getUserByID(uSvc)).Methods("GET")

	apiAuth.HandleFunc("/images/archives", archiveImage(jSvc)).Methods("POST")

	return &server{
		server: &http.Server{
			Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Handler: cors.New(cors.Options{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
				AllowCredentials: true,
			}).Handler(r),
		},
		done: make(chan struct{}),
	}
}
