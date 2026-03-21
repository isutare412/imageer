package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/isutare412/imageer/internal/processor/config"
	"github.com/isutare412/imageer/internal/processor/image"
	"github.com/isutare412/imageer/internal/processor/kafka"
	"github.com/isutare412/imageer/internal/processor/s3"
	imagesvc "github.com/isutare412/imageer/internal/processor/service/image"
	"github.com/isutare412/imageer/internal/processor/web"
)

type application struct {
	webServer     *web.Server
	kafkaClient   *kafka.Client
	kafkaConsumer *kafka.Consumer
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

	slog.Info("Create Kafka client")
	kafkaClient, err := kafka.NewClient(cfg.ToKafkaClientConfig())
	if err != nil {
		return nil, fmt.Errorf("creating kafka client: %w", err)
	}

	slog.Info("Create Kafka image process result queue")
	imageProcessResultQueue := kafka.NewImageProcessResultQueue(
		cfg.ToKafkaImageProcessResultQueueConfig(), kafkaClient)

	slog.Info("Create image service")
	imageService := imagesvc.NewService(imageProcessor, objectStorage, imageProcessResultQueue)

	slog.Info("Create Kafka image process request handler")
	imageProcessRequestHandler := kafka.NewImageProcessRequestHandler(
		cfg.ToKafkaImageProcessRequestHandlerConfig(), imageService)

	slog.Info("Create Kafka consumer")
	kafkaConsumer := kafka.NewConsumer(kafkaClient, map[string]kafka.Handler{
		cfg.Kafka.Topics.ImageProcessRequest.Topic:      imageProcessRequestHandler,
		cfg.Kafka.Topics.ImageProcessRequest.RetryTopic: imageProcessRequestHandler,
	})
	imageProcessRequestHandler.SetConsumer(kafkaConsumer)

	slog.Info("Create web server")
	webServer := web.NewServer(cfg.ToWebServerConfig())

	return &application{
		webServer:     webServer,
		kafkaClient:   kafkaClient,
		kafkaConsumer: kafkaConsumer,
	}, nil
}

func (a *application) initialize() error {
	return nil
}

func (a *application) run() {
	slog.Info("Run Kafka consumer")
	a.kafkaConsumer.Run()

	slog.Info("Run web server")
	webServerErrs := a.webServer.Run()

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

	slog.Info("Shutdown Kafka consumer")
	a.kafkaConsumer.Shutdown()

	slog.Info("Shutdown Kafka client")
	a.kafkaClient.Shutdown()
}

func logDuration(operation string) func() {
	start := time.Now()
	slog.Info(operation + " started")
	return func() {
		slog.Info(operation+" completed", "duration", time.Since(start))
	}
}
