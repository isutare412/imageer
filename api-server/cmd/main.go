package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/isutare412/imageer/api-server/pkg/adapter/http"
	"github.com/isutare412/imageer/api-server/pkg/adapter/mq"
	"github.com/isutare412/imageer/api-server/pkg/adapter/repository"
	"github.com/isutare412/imageer/api-server/pkg/config"
	"github.com/isutare412/imageer/api-server/pkg/core/encrypt"
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
		log.Fatalf("Need environment variable: %s", cfgEnvStr)
	}
	cfg, err := readConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	setLogger(cfg.Server.Mode)

	rootCtx, cancel := context.WithCancel(context.Background())

	redisMQ, err := mq.NewRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to create Redis MQ: %v", err)
	}
	log.Infof("Created Redis MQ on %v", cfg.Redis.Addrs)

	mysqlRepo, err := repository.NewMySQL(&cfg.MySQL)
	if err != nil {
		log.Fatalf("Failed to create MySQL repository: %v", err)
	}
	log.Infof("Created MySQL repository on %v", cfg.MySQL.Address)

	ecrSvc := encrypt.NewService()
	log.Info("Created encrypt service")

	uSvc := user.NewService(mysqlRepo, ecrSvc)
	log.Info("Created user service")

	jSvc := job.NewService(redisMQ)
	log.Info("Created job service")

	server := http.NewServer(&cfg.Server.Http, jSvc, uSvc)
	log.Info("Created HTTP server")

	// Start services
	sErr := server.Start(rootCtx)

	// Wait for signal or error
	sig := make(chan os.Signal, 3)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case e := <-sErr:
		log.Errorf("Got error from http server: %v", e)
	case s := <-sig:
		log.Infof("Caught signal: %s", s.String())
	}

	// Wait for graceful shutdown
	cancel()
	<-server.Done()
}

func readConfig(path string) (*config.Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("on readConfig: %w", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("on readConfig: %w", err)
	}
	return &cfg, nil
}

func setLogger(mode string) {
	if mode != "production" {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.TraceLevel)
	}
}
