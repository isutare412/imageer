package auth

import (
	"net/http"
	"time"
)

func (s *Service) createOIDCStateCookie(state string) *http.Cookie {
	return &http.Cookie{
		Name:     s.cfg.StateCookieName,
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(s.cfg.StateCookieTTL),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (s *Service) deleteOIDCStateCookie() *http.Cookie {
	return &http.Cookie{
		Name:     s.cfg.StateCookieName,
		Path:     "/",
		MaxAge:   -1, // Delete cookie immediately
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (s *Service) createUserCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     s.cfg.UserCookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(s.cfg.UserCookieTTL).Add(-time.Minute), // Delete cookie a bit earlier than token expiry
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (s *Service) deleteUserCookie() *http.Cookie {
	return &http.Cookie{
		Name:     s.cfg.UserCookieName,
		Path:     "/",
		MaxAge:   -1, // Delete cookie immediately
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}
