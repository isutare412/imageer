package log

import (
	"log/slog"
)

//go:generate go tool enumer -type=Level -trimprefix Level -output level_enum.go -transform lower -text
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

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
