package webv2

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	muxhandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/isutare412/imageer/internal/gateway/metric"
	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/internal/gateway/webv2/auth"
	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
	"github.com/isutare412/imageer/internal/gateway/webv2/handlers"
	"github.com/isutare412/imageer/internal/gateway/webv2/middleware"
)

//go:embed openapi.yaml openapi.html
var staticContents embed.FS

type Server struct {
	cfg    Config
	server *http.Server
}

func NewServer(
	cfg Config,
	healthCheckers []port.HealthChecker,
	authSvc port.AuthService,
	serviceAccountSvc port.ServiceAccountService,
	projectSvc port.ProjectService,
	userSvc port.UserService,
	imageSvc port.ImageService,
) (*Server, error) {
	handler := handlers.NewHandler(healthCheckers, authSvc, serviceAccountSvc, projectSvc, userSvc,
		imageSvc)

	authenticator := auth.NewAuthenticator(cfg.APIKeyHeader, cfg.UserCookieName,
		cfg.TokenRefreshThreshold, authSvc, serviceAccountSvc)

	authorizer := auth.NewAuthorizer(serviceAccountSvc, projectSvc, imageSvc)

	baseMiddlewares := []mux.MiddlewareFunc{
		middleware.ProxyHeaders,
		middleware.WithLogAttrContext,
		middleware.WithContextBag,
		middleware.WithRequestID,
		middleware.WithResponseRecord,
		middleware.AccessLog,
		middleware.ObserveMetrics,
		middleware.RecoverPanic,
		muxhandlers.CORS(cfg.CORS.buildCORSOptions()...),
	}

	apiMiddlewares := append(baseMiddlewares,
		authenticator.Authenticate,
		authorizer.Authorize,
		middleware.WithOpenAPIValidator())

	r := mux.NewRouter()

	r.HandleFunc("/healthz/live", handler.Liveness).Methods("GET")
	r.HandleFunc("/healthz/ready", handler.Readiness).Methods("GET")

	if cfg.ShowMetrics {
		r.Handle("/metrics", promhttp.HandlerFor(metric.Gatherer(),
			promhttp.HandlerOpts{})).
			Methods("GET")
	}

	if cfg.ShowOpenAPIDocs {
		docsRouter := r.PathPrefix("/docs").Subrouter()
		docsRouter.Use(baseMiddlewares...)

		docsRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/docs/openapi.html", http.StatusMovedPermanently)
		}).Methods("GET")

		docsRouter.PathPrefix("/").
			Handler(http.StripPrefix("/docs/",
				http.FileServer(http.FS(staticContents)))).
			Methods("GET")
	}

	// API routes - ALL middleware
	apiRouter := r.PathPrefix("/").Subrouter()
	apiRouter.Use(apiMiddlewares...)
	gen.HandlerWithOptions(handler, gen.GorillaServerOptions{
		BaseRouter:       apiRouter,
		ErrorHandlerFunc: gen.RespondError,
	})

	if err := logRegisteredRoutes(r); err != nil {
		return nil, fmt.Errorf("logging registerred routes: %w", err)
	}

	return &Server{
		cfg: cfg,
		server: &http.Server{
			Addr:              fmt.Sprintf(":%d", cfg.Port),
			Handler:           r,
			WriteTimeout:      cfg.WriteTimeout,
			ReadTimeout:       cfg.ReadTimeout,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		},
	}, nil
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

func logRegisteredRoutes(r *mux.Router) error {
	return r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		methods, err := route.GetMethods()
		if err != nil {
			// err is returned when the route has no specific method
			methods = []string{"ANY"}
		}

		for _, method := range methods {
			slog.Debug("Registered route", "method", method, "path", path)
		}

		return nil
	})
}
