package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/isutare412/imageer/image-processor/pkg/adapter/mq"
	"github.com/isutare412/imageer/image-processor/pkg/config"
	"github.com/isutare412/imageer/image-processor/pkg/core/job"
)

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

	config.SetMode(config.Mode(cfg.Mode))
	if config.IsDevelopmentMode() {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.TraceLevel)
	}

	rootCtx, cancel := context.WithCancel(context.Background())

	redisMQ, err := mq.NewRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("failed to create RedisMq: %v", err)
	}
	log.Infof("created redis MQ on %v", cfg.Redis.Addrs)

	pSvc, err := job.NewService(&cfg.Job, redisMQ)
	if err != nil {
		log.Fatalf("failed to create processor service: %v", err)
	}
	log.Info("created processor service")

	// Start services
	pSvc.Start(rootCtx)

	// Wait for signals
	sig := make(chan os.Signal, 3)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sig
	log.Infof("caught signal: %s", s.String())

	// Wait for graceful shutdown
	cancel()
	<-pSvc.Done()
}

func readConfig(path string) (*config.Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("on reading config: %w", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("on unmarshaling config: %w", err)
	}
	return &cfg, nil
}
