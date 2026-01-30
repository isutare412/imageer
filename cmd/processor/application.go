package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/isutare412/imageer/internal/processor/config"
	"github.com/isutare412/imageer/internal/processor/image"
	"github.com/isutare412/imageer/internal/processor/s3"
	imagesvc "github.com/isutare412/imageer/internal/processor/service/image"
	"github.com/isutare412/imageer/internal/processor/valkey"
	"github.com/isutare412/imageer/internal/processor/web"
)

type application struct {
	webServer                  *web.Server
	valkeyClient               *valkey.Client
	imageProcessRequestHandler *valkey.ImageProcessRequestHandler
}

func newApplication(cfg config.Config) (*application, error) {
	defer logDuration("Application creation")()

	slog.Info("Create image processor")
	imageProcessor := image.NewProcessor()

	slog.Info("Create S3 object storage")
	objectStorage, err := s3.NewObjectStorage(cfg.ToS3ObjectStorageConfig())
	if err != nil {
		return nil, fmt.Errorf("creating s3 object storage: %w", err)
	}

	slog.Info("Create valkey client")
	valkeyClient, err := valkey.NewClient(cfg.ToValkeyClientConfig())
	if err != nil {
		return nil, fmt.Errorf("creating valkey client: %w", err)
	}

	slog.Info("Create valkey image process result queue")
	imageProcessResultQueue := valkey.NewImageProcessResultQueue(
		cfg.ToValkeyImageProcessResultQueueConfig(), valkeyClient)

	slog.Info("Create image service")
	imageService := imagesvc.NewService(imageProcessor, objectStorage, imageProcessResultQueue)

	slog.Info("Create valkey image process request handler")
	imageProcessRequestHandler := valkey.NewImageProcessRequestHandler(
		cfg.ToValkeyImageProcessRequestHandlerConfig(), valkeyClient, imageService)

	slog.Info("Create web server")
	webServer := web.NewServer(cfg.ToWebServerConfig())

	return &application{
		webServer:                  webServer,
		valkeyClient:               valkeyClient,
		imageProcessRequestHandler: imageProcessRequestHandler,
	}, nil
}

func (a *application) initialize() error {
	defer logDuration("Application initialization")()

	ctx, cancelTimeout := context.WithTimeout(context.Background(), time.Minute)
	defer cancelTimeout()

	ctx, cancelSignal := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancelSignal()

	slog.Info("Initialize image process request handler")
	if err := a.imageProcessRequestHandler.Initialize(ctx); err != nil {
		return fmt.Errorf("initializing image process request handler: %w", err)
	}

	return nil
}

func (a *application) run() {
	slog.Info("Run web server")
	webServerErrs := a.webServer.Run()

	slog.Info("Run image process request handler")
	a.imageProcessRequestHandler.Run()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case sig := <-signals:
		slog.Info("Received signal; shutdown application", "signal", sig)
	case err := <-webServerErrs:
		slog.Error("Web server error; shutdown application", "error", err)
	}
}

func (a *application) shutdown() {
	defer logDuration("Application shutdown")()

	slog.Info("Shutdown web server")
	if err := a.webServer.Shutdown(); err != nil {
		slog.Error("Failed to shutdown web server", "error", err)
	}

	slog.Info("Shutdown image process request handler")
	a.imageProcessRequestHandler.Shutdown()

	slog.Info("Shutdown valkey client")
	a.valkeyClient.Shutdown()
}

func logDuration(operation string) func() {
	start := time.Now()
	slog.Info(operation + " started")
	return func() {
		slog.Info(operation+" completed", "duration", time.Since(start))
	}
}
