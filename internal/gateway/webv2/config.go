package webv2

import (
	"time"

	"github.com/gorilla/handlers"
)

type Config struct {
	Port                  int
	ShowMetrics           bool
	ShowOpenAPIDocs       bool
	APIKeyHeader          string
	UserCookieName        string
	TokenRefreshThreshold time.Duration
	WriteTimeout          time.Duration
	ReadTimeout           time.Duration
	ReadHeaderTimeout     time.Duration
	CORS                  CORSConfig
}

type CORSConfig struct {
	AllowOrigins     []string
	AllowHeaders     []string
	AllowMethods     []string
	AllowCredentials bool
	MaxAge           time.Duration
}

func (c CORSConfig) buildCORSOptions() []handlers.CORSOption {
	opts := []handlers.CORSOption{
		handlers.AllowedOrigins(c.AllowOrigins),
		handlers.AllowedHeaders(c.AllowHeaders),
		handlers.AllowedMethods(c.AllowMethods),
		handlers.MaxAge(int(c.MaxAge.Seconds())),
	}
	if c.AllowCredentials {
		opts = append(opts, handlers.AllowCredentials())
	}
	return opts
}
