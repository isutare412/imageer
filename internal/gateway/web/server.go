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
)

//go:embed openapi.yaml openapi.html
var staticContents embed.FS

type Server struct {
	cfg     Config
	engine  *echo.Echo
	handler *handler
}

func NewServer(cfg Config, authSvc port.AuthService) *Server {
	handler := newHandler(authSvc)

	e := echo.New()
	e.HidePort = true
	e.HideBanner = !cfg.ShowBanner
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	e.Use(
		withLogAttrContext,
		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			RequestIDHandler: attachRequestID,
		}),
		accessLog,
		respondError,
		recoverPanic,
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
