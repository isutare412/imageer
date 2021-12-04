package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/isutare412/imageer/image-processor/pkg/adapter/mq"
	"github.com/isutare412/imageer/image-processor/pkg/config"
)

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

	_, err = mq.NewRedisMq(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to create RedisMq: %v", err)
	}
	log.Info("Created Redis MQ")
}

func readConfig(path string) (*config.Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("on reading config: %v", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("on unmarshaling config: %v", err)
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
