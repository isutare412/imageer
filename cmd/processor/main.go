package main

import (
	"flag"
	"log/slog"

	pcconfig "github.com/isutare412/imageer/internal/processor/config"
	"github.com/isutare412/imageer/pkg/config"
	"github.com/isutare412/imageer/pkg/log"
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
