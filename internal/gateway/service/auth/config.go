package auth

import "time"

type ServiceConfig struct {
	StateCookieName string
	StateCookieTTL  time.Duration
	UserCookieName  string
	UserCookieTTL   time.Duration
}
