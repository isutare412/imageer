package postgres

import sloggorm "github.com/orandin/slog-gorm"

type ClientConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	Database    string
	TraceLog    bool
	UseInMemory bool
}

func (c *ClientConfig) buildSlogGORMOption() []sloggorm.Option {
	var opts []sloggorm.Option
	if c.TraceLog {
		opts = append(opts, sloggorm.WithTraceAll())
	}
	return opts
}
