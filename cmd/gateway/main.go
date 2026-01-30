package main

import (
	"flag"
	"log/slog"

	gwconfig "github.com/isutare412/imageer/internal/gateway/config"
	"github.com/isutare412/imageer/internal/gateway/metric"
	"github.com/isutare412/imageer/pkg/config"
	"github.com/isutare412/imageer/pkg/log"
)

var cfgPath = flag.String("configs", ".", "Path to config directory")

func init() {
	flag.Parse()
}

func main() {
	cfg, err := config.LoadValidated[gwconfig.Config]("GATEWAY", *cfgPath)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return
	}

	log.Init(cfg.ToLogConfig())
	if cfg.Metrics.Enabled {
		metric.Init()
	}

	slog.Debug("Loaded config", "config", cfg)

	app, err := newApplication(cfg)
	if err != nil {
		slog.Error("Failed to create application", "error", err)
		return
	}

	if err := app.initialize(); err != nil {
		slog.Error("Failed to initialize application", "error", err)
		return
	}

	app.run()
	app.shutdown()
}
