package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/isutare412/imageer/api-server/pkg/adapter/http"
	"github.com/isutare412/imageer/api-server/pkg/adapter/mq"
	"github.com/isutare412/imageer/api-server/pkg/adapter/repo"
	"github.com/isutare412/imageer/api-server/pkg/config"
	"github.com/isutare412/imageer/api-server/pkg/core/auth"
	"github.com/isutare412/imageer/api-server/pkg/core/job"
	"github.com/isutare412/imageer/api-server/pkg/core/user"
)

// @title Imageer Endpoint API
// @version 0.1
// @description Endpoint API for image processing service.
func main() {
	const cfgEnvStr = "IMAGEER_CONFIG"

	cfgPath := os.Getenv(cfgEnvStr)
	if cfgPath == "" {
		log.Fatalf("need environment variable: %s", cfgEnvStr)
	}
	cfg, err := readConfig(cfgPath)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	setLogger(cfg.Server.Mode)

	rootCtx, cancel := context.WithCancel(context.Background())

	redisMQ, err := mq.NewRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("failed to create Redis MQ: %v", err)
	}
	log.Infof("created Redis MQ on %v", cfg.Redis.Addrs)

	mysqlRepo, err := repo.NewMySQL(&cfg.MySQL)
	if err != nil {
		log.Fatalf("failed to create MySQL repository: %v", err)
	}
	log.Infof("created MySQL repository on %v", cfg.MySQL.Address)

	s3Repo, err := repo.NewS3(&cfg.S3)
	if err != nil {
		log.Fatalf("failed to create S3 repository: %v", err)
	}
	log.Infof("created S3 repository on %v", cfg.S3.Address)

	authSvc, err := auth.NewService(&cfg.Auth)
	if err != nil {
		log.Fatalf("failed to create auth service: %v", err)
	}
	log.Info("created auth service")

	uSvc := user.NewService(mysqlRepo, authSvc)
	log.Info("created user service")

	jSvc := job.NewService(&cfg.Server.Job, redisMQ, s3Repo)
	log.Info("created job service")

	server := http.NewServer(&cfg.Server.Http, jSvc, uSvc, authSvc)
	log.Info("created HTTP server")

	// Start services
	sErr := server.Start(rootCtx)

	// Wait for signal or error
	sig := make(chan os.Signal, 3)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case e := <-sErr:
		log.Errorf("got error from http server: %v", e)
	case s := <-sig:
		log.Infof("caught signal[%s]", s.String())
	}

	// Wait for graceful shutdown
	cancel()
	<-server.Done()
}

func readConfig(path string) (*config.Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func setLogger(mode string) {
	if mode == "development" {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	}
}
