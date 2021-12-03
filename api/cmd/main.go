package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/isutare412/imageer/api/internal/adapter/http"
	"github.com/isutare412/imageer/api/internal/config"
)

// @title Imageer Endpoint API
// @version 0.1
// @description Endpoint API for image processing service.

// @host localhost:8080
// @BasePath /api/v1
func main() {
	const ENV_CFG = "IMAGEER_CONFIG"

	cPath := os.Getenv(ENV_CFG)
	if cPath == "" {
		log.Fatalf("Need environment variable: %s", ENV_CFG)
	}
	cfg, err := readConfig(cPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	setLogger(cfg.Server.Mode)

	server := http.New(&cfg.Server.Http)
	sErr := server.Start()

	sig := make(chan os.Signal, 3)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case e := <-sErr:
		log.Errorf("Got error from http server: %v", e)
	case s := <-sig:
		log.Infof("Caught signal: %s", s.String())
	}

	server.Shutdown()
}

func readConfig(path string) (*config.Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("on readConfig: %v", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("on readConfig: %v", err)
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
