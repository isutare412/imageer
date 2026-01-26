package webv2

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
	"github.com/isutare412/imageer/internal/gateway/webv2/immigration"
	"github.com/isutare412/imageer/internal/gateway/webv2/middleware"
)

//go:embed openapi.yaml openapi.html
var staticContents embed.FS

type Server struct {
	cfg    Config
	server *http.Server
}

func NewServer(
	cfg Config, authSvc port.AuthService, serviceAccountSvc port.ServiceAccountService,
	projectSvc port.ProjectService, userSvc port.UserService, imageSvc port.ImageService,
) (*Server, error) {
	handler := newHandler(authSvc, serviceAccountSvc, projectSvc, userSvc, imageSvc)

	passportIssuer := immigration.NewPassportIssuer(cfg.APIKeyHeader, cfg.UserCookieName,
		authSvc, serviceAccountSvc)

	immigration := immigration.New(serviceAccountSvc, projectSvc, imageSvc)

	middlewares := []mux.MiddlewareFunc{
		handlers.ProxyHeaders,
		middleware.WithLogAttrContext,
		middleware.WithContextBag,
		middleware.WithRequestID,
		middleware.WithResponseMetrics,
		middleware.AccessLog,
		middleware.RecoverPanic,
		handlers.CORS(cfg.CORS.buildCORSOptions()...),
		middleware.WithOpenAPIValidator(),
		passportIssuer.IssuePassport,
		immigration.Immigrate,
	}

	r := mux.NewRouter()
	r.Use(middlewares...)

	gen.HandlerWithOptions(handler, gen.GorillaServerOptions{
		BaseRouter:       r,
		ErrorHandlerFunc: gen.RespondError,
	})

	if cfg.ShowOpenAPIDocs {
		r.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/docs/openapi.html", http.StatusMovedPermanently)
		}).Methods("GET")

		r.PathPrefix("/docs/").
			Handler(http.StripPrefix("/docs/",
				http.FileServer(http.FS(staticContents)))).
			Methods("GET")
	}

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
			slog.Debug("Regiestered route", "method", method, "path", path)
		}

		return nil
	})
}
