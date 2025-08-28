package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/isutare412/imageer/internal/gateway/web"
	"github.com/isutare412/imageer/pkg/log"
)

func main() {
	log.Init(log.Config{
		Format:    log.FormatPretty,
		Level:     log.LevelDebug,
		AddSource: true,
	})

	webServer := web.NewServer(web.Config{
		Port:            8080,
		ShowBanner:      true,
		ShowOpenAPIDocs: true,
	})

	webServerErrs := webServer.Run()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case sig := <-signals:
		slog.Info("Received signal, shutting down", "signal", sig)
	case err := <-webServerErrs:
		slog.Error("Web server error", "error", err)
	}

	if err := webServer.Shutdown(); err != nil {
		slog.Error("Failed to shutdown web server", "error", err)
	}
}
