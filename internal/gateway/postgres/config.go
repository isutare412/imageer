package postgres

import (
	"fmt"

	sloggorm "github.com/orandin/slog-gorm"
)

type ClientConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	Database    string
	SSLMode     SSLMode
	TraceLog    bool
	UseInMemory bool
}

// SSLMode is the libpq sslmode. Certificate verification modes are not
// supported because the client does not load a CA certificate.
type SSLMode string

const (
	SSLModeDisable SSLMode = "disable"
	SSLModeAllow   SSLMode = "allow"
	SSLModePrefer  SSLMode = "prefer"
	SSLModeRequire SSLMode = "require"
)

// Validate accepts an empty value, which means disable.
func (m SSLMode) Validate() error {
	switch m {
	case "", SSLModeDisable, SSLModeAllow, SSLModePrefer, SSLModeRequire:
	default:
		return fmt.Errorf("unexpected ssl mode %q", m)
	}
	return nil
}

// Value returns the sslmode for the DSN. Empty value falls back to disable.
func (m SSLMode) Value() string {
	if m == "" {
		return string(SSLModeDisable)
	}
	return string(m)
}

func (c *ClientConfig) buildSlogGORMOption() []sloggorm.Option {
	var opts []sloggorm.Option
	if c.TraceLog {
		opts = append(opts, sloggorm.WithTraceAll())
	}
	return opts
}
