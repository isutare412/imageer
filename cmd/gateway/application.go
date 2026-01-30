package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/isutare412/imageer/internal/gateway/config"
	"github.com/isutare412/imageer/internal/gateway/crypt"
	"github.com/isutare412/imageer/internal/gateway/jwt"
	"github.com/isutare412/imageer/internal/gateway/kubernetes"
	"github.com/isutare412/imageer/internal/gateway/oidc"
	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/internal/gateway/postgres"
	"github.com/isutare412/imageer/internal/gateway/s3"
	"github.com/isutare412/imageer/internal/gateway/service/auth"
	"github.com/isutare412/imageer/internal/gateway/service/image"
	"github.com/isutare412/imageer/internal/gateway/service/project"
	"github.com/isutare412/imageer/internal/gateway/service/serviceaccount"
	"github.com/isutare412/imageer/internal/gateway/service/user"
	"github.com/isutare412/imageer/internal/gateway/sqs"
	"github.com/isutare412/imageer/internal/gateway/valkey"
	"github.com/isutare412/imageer/internal/gateway/webv2"
)

type application struct {
	webServer                   *webv2.Server
	imageUploadListener         *sqs.ImageUploadListener
	postgresClient              *postgres.Client
	valkeyClient                *valkey.Client
	imageProcessResultHandler   *valkey.ImageProcessResultHandler
	imageS3DeleteRequestHandler *valkey.ImageS3DeleteRequestHandler
	leaderElector               leaderElector

	cfg config.Config
}

func newApplication(cfg config.Config) (*application, error) {
	defer logDuration("Application creation")()

	slog.Info("Create OIDC Google client")
	oidcProvider, err := oidc.NewGoogleClient(cfg.ToOIDCGoogleClientConfig())
	if err != nil {
		return nil, fmt.Errorf("creating oidc google client: %w", err)
	}

	slog.Info("Create AES crypter")
	aesCrypter, err := crypt.NewAESCrypter(cfg.ToAESCrypterConfig())
	if err != nil {
		return nil, fmt.Errorf("creating aes crypter: %w", err)
	}

	slog.Info("Create JWT signer")
	jwtSigner, err := jwt.NewSigner(cfg.ToJWTSignerConfig())
	if err != nil {
		return nil, fmt.Errorf("creating jwt signer: %w", err)
	}

	slog.Info("Create JWT verifier")
	jwtVerifier, err := jwt.NewVerifier(cfg.ToJWTVerifierConfig())
	if err != nil {
		return nil, fmt.Errorf("creating jwt verifier: %w", err)
	}

	slog.Info("Create s3 presigner")
	s3Presigner, err := s3.NewPresigner(cfg.ToS3PresignerConfig())
	if err != nil {
		return nil, fmt.Errorf("creating s3 presigner: %w", err)
	}

	slog.Info("Create s3 object storage")
	s3ObjectStorage, err := s3.NewObjectStorage(cfg.ToS3ObjectStorageConfig())
	if err != nil {
		return nil, fmt.Errorf("creating s3 object storage: %w", err)
	}

	slog.Info("Create repository client")
	postgresClient, err := postgres.NewClient(cfg.ToRepositoryClientConfig())
	if err != nil {
		return nil, fmt.Errorf("creating repository client: %w", err)
	}

	slog.Info("Create transactioner")
	transactioner := postgres.NewTransactioner(postgresClient)

	slog.Info("Create user repository")
	userRepo := postgres.NewUserRepository(postgresClient)

	slog.Info("Create service account repository")
	serviceAccountRepo := postgres.NewServiceAccountRepository(postgresClient)

	slog.Info("Create project repository")
	projectRepo := postgres.NewProjectRepository(postgresClient)

	slog.Info("Create image repository")
	imageRepo := postgres.NewImageRepository(postgresClient)

	slog.Info("Create image variant repository")
	imageVarRepo := postgres.NewImageVariantRepository(postgresClient)

	slog.Info("Create image processing log repository")
	imageProcLogRepo := postgres.NewImageProcessingLogRepository(postgresClient)

	slog.Info("Create preset repository")
	presetRepo := postgres.NewPresetRepository(postgresClient)

	slog.Info("Create valkey client")
	valkeyClient, err := valkey.NewClient(cfg.ToValkeyClientConfig())
	if err != nil {
		return nil, fmt.Errorf("creating valkey client: %w", err)
	}

	slog.Info("Create valkey image process request queue")
	imageProcRequestQueue := valkey.NewImageProcessRequestQueue(
		cfg.ToValkeyImageProcessRequestQueueConfig(), valkeyClient)

	slog.Info("Create valkey image S3 delete request queue")
	imageS3DeleteRequestQueue := valkey.NewImageS3DeleteRequestQueue(
		cfg.ToValkeyImageS3DeleteRequestQueueConfig(), valkeyClient)

	slog.Info("Create valkey image notification publisher")
	imageNotificationPublisher := valkey.NewImageNotificationPublisher(
		cfg.ToValkeyImageNotificationPublisherConfig(), valkeyClient)

	slog.Info("Create valkey image upload done subscriber")
	imageUploadDoneSubscriber := valkey.NewImageUploadDoneSubscriber(
		cfg.ToValkeyImageUploadDoneSubscriberConfig(), valkeyClient)

	slog.Info("Create valkey image process done subscriber")
	imageProcDoneSubscriber := valkey.NewImageProcessDoneSubscriber(
		cfg.ToValkeyImageProcessDoneSubscriberConfig(), valkeyClient)

	slog.Info("Create auth service")
	authSvc := auth.NewService(cfg.ToAuthServiceConfig(), oidcProvider,
		aesCrypter, jwtSigner, jwtVerifier, userRepo)

	slog.Info("Create service account service")
	serviceAccountSvc := serviceaccount.NewService(serviceAccountRepo)

	slog.Info("Create project service")
	projectSvc := project.NewService(projectRepo)

	slog.Info("Create user service")
	userSvc := user.NewService(userRepo)

	slog.Info("Create image service")
	imageSvc := image.NewService(cfg.ToImageServiceConfig(), s3Presigner, s3ObjectStorage,
		transactioner, imageRepo, imageVarRepo, imageProcLogRepo, presetRepo, imageProcRequestQueue,
		imageNotificationPublisher, imageUploadDoneSubscriber, imageProcDoneSubscriber,
		imageS3DeleteRequestQueue)

	slog.Info("Create web server")
	healthCheckers := []port.HealthChecker{postgresClient, valkeyClient}
	webServer, err := webv2.NewServer(cfg.ToWebV2Config(), healthCheckers, authSvc,
		serviceAccountSvc, projectSvc, userSvc, imageSvc)
	if err != nil {
		return nil, fmt.Errorf("creating web server: %w", err)
	}

	slog.Info("Create SQS image upload listener")
	imageUploadListener, err := sqs.NewImageUploadListener(cfg.ToSQSImageUploadListenerConfig(),
		imageSvc)
	if err != nil {
		return nil, fmt.Errorf("creating SQS image upload listener: %w", err)
	}

	slog.Info("Create valkey image process result handler")
	imageProcessResultHandler := valkey.NewImageProcessResultHandler(
		cfg.ToValkeyImageProcessResultHandlerConfig(), valkeyClient, imageSvc)

	slog.Info("Create valkey image S3 delete request handler")
	imageS3DeleteRequestHandler := valkey.NewImageS3DeleteRequestHandler(
		cfg.ToValkeyImageS3DeleteRequestHandlerConfig(), valkeyClient, imageSvc)

	slog.Info("Create image closer")
	imageCloser := image.NewCloser(cfg.ToImageCloserConfig(), transactioner, imageRepo,
		imageVarRepo)

	handlers := []port.LeaderHandler{imageCloser}

	var elector leaderElector
	if cfg.Kubernetes.Enabled {
		slog.Info("Create kubernetes client")
		k8sClient, err := kubernetes.NewClient(cfg.ToKubernetesClientConfig())
		if err != nil {
			return nil, fmt.Errorf("creating kubernetes client: %w", err)
		}

		slog.Info("Create kubernetes leader elector")
		elector, err = kubernetes.NewLeaderElector(cfg.ToKubernetesLeaderElectorConfig(),
			k8sClient, handlers)
		if err != nil {
			return nil, fmt.Errorf("creating kubernetes leader elector: %w", err)
		}
	} else {
		slog.Info("Create standalone elector")
		elector = kubernetes.NewStandaloneElector(handlers)
	}

	return &application{
		webServer:                   webServer,
		imageUploadListener:         imageUploadListener,
		postgresClient:              postgresClient,
		valkeyClient:                valkeyClient,
		imageProcessResultHandler:   imageProcessResultHandler,
		imageS3DeleteRequestHandler: imageS3DeleteRequestHandler,
		leaderElector:               elector,
		cfg:                         cfg,
	}, nil
}

func (a *application) initialize() error {
	defer logDuration("Application initialization")()

	ctx, cancelTimeout := context.WithTimeout(context.Background(), time.Minute)
	defer cancelTimeout()

	ctx, cancelSignal := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancelSignal()

	slog.Info("Migrate database schemas")
	if err := a.postgresClient.MigrateSchemas(ctx); err != nil {
		return fmt.Errorf("migrating database schemas: %w", err)
	}

	slog.Info("Initialize image process result handler")
	if err := a.imageProcessResultHandler.Initialize(ctx); err != nil {
		return fmt.Errorf("initializing image process result handler: %w", err)
	}

	slog.Info("Initialize image S3 delete request handler")
	if err := a.imageS3DeleteRequestHandler.Initialize(ctx); err != nil {
		return fmt.Errorf("initializing image S3 delete request handler: %w", err)
	}

	return nil
}

func (a *application) run() {
	slog.Info("Run image upload listener")
	a.imageUploadListener.Run()

	slog.Info("Run image process result handler")
	a.imageProcessResultHandler.Run()

	slog.Info("Run image S3 delete request handler")
	a.imageS3DeleteRequestHandler.Run()

	slog.Info("Run kubernetes leader elector")
	a.leaderElector.Run()

	slog.Info("Run web server")
	webServerErrs := a.webServer.Run()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case sig := <-signals:
		slog.Info("Received signal; shutdown application", "signal", sig)
	case err := <-webServerErrs:
		slog.Error("Detected web server error; shutdown application", "error", err)
	}
}

func (a *application) shutdown() {
	defer logDuration("Application shutdown")()

	slog.Info("Shutdown web server")
	if err := a.webServer.Shutdown(); err != nil {
		slog.Error("Failed to shutdown web server", "error", err)
	}

	slog.Info("Shutdown kubernetes leader elector")
	a.leaderElector.Shutdown()

	slog.Info("Shutdown image upload listener")
	a.imageUploadListener.Shutdown()

	slog.Info("Shutdown image process result handler")
	a.imageProcessResultHandler.Shutdown()

	slog.Info("Shutdown image S3 delete request handler")
	a.imageS3DeleteRequestHandler.Shutdown()

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
