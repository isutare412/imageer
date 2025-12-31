package log

import (
	"log/slog"

	"github.com/isutare412/imageer/pkg/apperr"
)

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

func (l Level) Validate() error {
	switch l {
	case LevelDebug:
	case LevelInfo:
	case LevelWarn:
	case LevelError:
	default:
		return apperr.NewError(apperr.CodeBadRequest).WithSummary("Unexpected log level %q", l)
	}
	return nil
}

func (l Level) SlogLevel() slog.Level {
	sl := slog.LevelInfo
	switch l {
	case LevelDebug:
		sl = slog.LevelDebug
	case LevelInfo:
		sl = slog.LevelInfo
	case LevelWarn:
		sl = slog.LevelWarn
	case LevelError:
		sl = slog.LevelError
	}
	return sl
}
