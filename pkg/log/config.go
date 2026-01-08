package log

import "log/slog"

type Config struct {
	Format      Format // Optional. Defaults to text.
	Level       Level  // Optional. Defaults to info.
	AddSource   bool   // Optional. Defaults to false.
	Component   string // Optional
	Environment string // Optional
	Version     string // Optional
}

func (c Config) ConstAttrs() []slog.Attr {
	var attrs []slog.Attr
	if c.Component != "" {
		attrs = append(attrs, slog.String("component", c.Component))
	}
	if c.Environment != "" {
		attrs = append(attrs, slog.String("environment", c.Environment))
	}
	if c.Version != "" {
		attrs = append(attrs, slog.String("version", c.Version))
	}
	return attrs
}
