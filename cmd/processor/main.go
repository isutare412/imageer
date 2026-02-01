package main

import (
	"flag"
	"log/slog"

	pcconfig "github.com/isutare412/imageer/internal/processor/config"
	"github.com/isutare412/imageer/internal/processor/metric"
	"github.com/isutare412/imageer/pkg/config"
	"github.com/isutare412/imageer/pkg/log"
	"github.com/isutare412/imageer/pkg/tracing"
)

var cfgPath = flag.String("configs", ".", "Path to config directory")

func init() {
	flag.Parse()
}

func main() {
	cfg, err := config.LoadValidated[pcconfig.Config]("PROCESSOR", *cfgPath)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return
	}

	log.Init(cfg.ToLogConfig())
	metric.Init()

	if err := tracing.Init(cfg.ToTracingConfig()); err != nil {
		slog.Error("Failed to initialize trace", "error", err)
		return
	}
	defer func() {
		if err := tracing.Shutdown(); err != nil {
			slog.Error("Failed to shutdown trace", "error", err)
		}
	}()

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
