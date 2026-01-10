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
	"github.com/isutare412/imageer/internal/gateway/oidc"
	"github.com/isutare412/imageer/internal/gateway/postgres"
	"github.com/isutare412/imageer/internal/gateway/s3"
	"github.com/isutare412/imageer/internal/gateway/service/auth"
	"github.com/isutare412/imageer/internal/gateway/service/image"
	"github.com/isutare412/imageer/internal/gateway/service/project"
	"github.com/isutare412/imageer/internal/gateway/service/serviceaccount"
	"github.com/isutare412/imageer/internal/gateway/sqs"
	"github.com/isutare412/imageer/internal/gateway/valkey"
	"github.com/isutare412/imageer/internal/gateway/web"
)

type application struct {
	webServer         *web.Server
	imageUploadLister *sqs.ImageUploadListener
	postgresClient    *postgres.Client
	valkeyClient      *valkey.Client

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

	slog.Info("Create preset repository")
	presetRepo := postgres.NewPresetRepository(postgresClient)

	slog.Info("Create valkey client")
	valkeyClient, err := valkey.NewClient(cfg.ToValkeyClientConfig())
	if err != nil {
		return nil, fmt.Errorf("creating valkey client: %w", err)
	}

	slog.Info("Create valkey image event queue")
	imageEventQueue := valkey.NewImageEventQueue(cfg.ToValkeyImageEventQueueConfig(), valkeyClient)

	slog.Info("Create auth service")
	authSvc := auth.NewService(cfg.ToAuthServiceConfig(), oidcProvider,
		aesCrypter, jwtSigner, jwtVerifier, userRepo)

	slog.Info("Create service account service")
	serviceAccountSvc := serviceaccount.NewService(serviceAccountRepo)

	slog.Info("Create project service")
	projectSvc := project.NewService(projectRepo)

	imageSvc := image.NewService(cfg.ToImageServiceConfig(), s3Presigner, transactioner, imageRepo,
		imageVarRepo, presetRepo, imageEventQueue)

	slog.Info("Create web server")
	webServer := web.NewServer(cfg.ToWebConfig(), authSvc, serviceAccountSvc, projectSvc, imageSvc)

	slog.Info("Create SQS image upload listener")
	imageUploadListener, err := sqs.NewImageUploadListener(cfg.ToSQSImageUploadListenerConfig(),
		imageSvc)
	if err != nil {
		return nil, fmt.Errorf("creating SQS image upload listener: %w", err)
	}

	return &application{
		webServer:         webServer,
		imageUploadLister: imageUploadListener,
		postgresClient:    postgresClient,
		valkeyClient:      valkeyClient,
		cfg:               cfg,
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

	return nil
}

func (a *application) run() {
	slog.Info("Run SQS image upload listener")
	a.imageUploadLister.Run()

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

	slog.Info("Shutdown SQS image upload listener")
	a.imageUploadLister.Shutdown()

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
