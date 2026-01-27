package web

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/internal/gateway/web/auth"
)

//go:embed openapi.yaml openapi.html
var staticContents embed.FS

type Server struct {
	cfg     Config
	engine  *echo.Echo
	handler *handler
}

func NewServer(
	cfg Config, authSvc port.AuthService, serviceAccountSvc port.ServiceAccountService,
	projectSvc port.ProjectService, userSvc port.UserService, imageSvc port.ImageService,
) *Server {
	handler := newHandler(authSvc, serviceAccountSvc, projectSvc, userSvc, imageSvc)

	authenticator := auth.NewAuthenticator(cfg.APIKeyHeader, cfg.UserCookieName,
		authSvc, serviceAccountSvc)

	authorizer := auth.NewAuthorizer(serviceAccountSvc, projectSvc, imageSvc)

	e := echo.New()
	e.HidePort = true
	e.HideBanner = !cfg.ShowBanner
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	e.Server.WriteTimeout = cfg.WriteTimeout
	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.ReadHeaderTimeout = cfg.ReadHeaderTimeout
	e.Use(
		withContextBag,
		withLogAttrContext,
		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			RequestIDHandler: attachRequestID,
		}),
		accessLog,
		respondError,
		recoverPanic,
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     cfg.CORS.AllowOrigins,
			AllowHeaders:     cfg.CORS.AllowHeaders,
			AllowMethods:     cfg.CORS.AllowMethods,
			AllowCredentials: cfg.CORS.AllowCredentials,
			MaxAge:           int(cfg.CORS.MaxAge.Seconds()),
		}),
		authenticator.Authenticate,
		authorizer.Authorize,
		openAPIValidator(),
	)

	RegisterHandlers(e, handler)

	if cfg.ShowOpenAPIDocs {
		e.FileFS("/docs/openapi.yaml", "openapi.yaml", &staticContents)
		e.FileFS("/docs", "openapi.html", &staticContents)
	}

	for _, r := range e.Routes() {
		slog.Debug("Registered route", "method", r.Method, "path", r.Path)
	}

	return &Server{
		cfg:     cfg,
		engine:  e,
		handler: handler,
	}
}

func (s *Server) Run() <-chan error {
	errs := make(chan error, 1)
	go func() {
		defer close(errs)

		slog.Info("Starting web server", "port", s.cfg.Port)
		if err := s.engine.Start(fmt.Sprintf(":%d", s.cfg.Port)); !errors.Is(err, http.ErrServerClosed) {
			errs <- fmt.Errorf("failed to start server: %w", err)
			return
		}
	}()

	return errs
}

func (s *Server) Shutdown() error {
	return s.engine.Shutdown(context.Background())
}
