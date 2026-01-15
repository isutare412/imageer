package web

import "time"

type Config struct {
	Port              int
	ShowBanner        bool
	ShowOpenAPIDocs   bool
	APIKeyHeader      string
	UserCookieName    string
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}
