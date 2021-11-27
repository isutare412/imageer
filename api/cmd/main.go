package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/isutare412/imageer/api/pkg/config"
)

const ENV_CFG = "IMAGEER_CONFIG"

func main() {
	cPath := os.Getenv(ENV_CFG)
	if cPath == "" {
		log.Fatalf("Need %q as environment variable", ENV_CFG)
	}
	cfg, err := readConfig(cPath)
	if err != nil {
		log.Fatalf("Failed to read config: %s", err)
	}

	setLogger(cfg.Server.Mode)

	log.Infof("http.host: %s", cfg.Server.Http.Host)
	log.Infof("http.port: %s", cfg.Server.Http.Port)
}

func readConfig(path string) (*config.Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("on readConfig: %s", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("on readConfig: %s", err)
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
