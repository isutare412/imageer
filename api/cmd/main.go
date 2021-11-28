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

	s := http.New(&cfg.Server.Http)
	sErrChan := s.Start()

	sigChan := make(chan os.Signal, 3)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case err := <-sErrChan:
		log.Errorf("Got error from http server: %v", err)
	case sig := <-sigChan:
		log.Infof("Caught signal: %s", sig.String())
	}

	s.Shutdown()
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
