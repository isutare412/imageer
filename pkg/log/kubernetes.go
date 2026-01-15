package log

import (
	"log/slog"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

func adaptKlog(handler slog.Handler) {
	klog.SetLogger(logr.FromSlogHandler(handler))
}
