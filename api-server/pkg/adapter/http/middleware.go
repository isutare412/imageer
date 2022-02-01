package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/isutare412/imageer/api-server/pkg/core/auth"
	"github.com/isutare412/imageer/api-server/pkg/core/user"
	log "github.com/sirupsen/logrus"
)

type responseLogger struct {
	http.ResponseWriter
	status int
	length int
}

func (l *responseLogger) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}

func (l *responseLogger) Write(b []byte) (int, error) {
	l.length += len(b)
	return l.ResponseWriter.Write(b)
}

func structAccessLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := responseLogger{ResponseWriter: w, status: http.StatusOK}
		h.ServeHTTP(&logger, r)

		log.WithFields(log.Fields{
			"addr":   r.RemoteAddr,
			"method": r.Method,
			"url":    r.URL.String(),
			"status": logger.status,
			"length": logger.length,
		}).Info("access")
	})
}

func plainAccessLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := responseLogger{ResponseWriter: w, status: http.StatusOK}
		h.ServeHTTP(&logger, r)

		log.Infof("%s - \"%s %s\" %d %d",
			r.RemoteAddr, r.Method, r.URL.String(), logger.status, logger.length)
	})
}

func checkSession(authSvc auth.Service) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string
			if cookie, err := r.Cookie("token"); err == nil {
				token = cookie.Value
			} else if authStr := r.Header.Get("Authorization"); authStr != "" {
				authSplit := strings.SplitN(authStr, "Bearer ", 2)
				if len(authSplit) < 2 {
					log.Warnf("failed to split authorization header")
					responseError(w, http.StatusUnauthorized, "invalid authorization header")
					return
				}
				token = authSplit[1]
			} else {
				log.Warnf("received request without auth")
				responseError(w, http.StatusUnauthorized, "need token from cookie or authorization header")
				return
			}

			sess, err := authSvc.VerifyToken(auth.Token(token))
			if errors.Is(err, auth.ErrTokenExpired) {
				log.Warnf("received expired token")
				responseError(w, http.StatusUnauthorized, "token expired")
				return
			} else if err != nil {
				log.Errorf("failed to verify token: %v", err)
				responseError(w, http.StatusUnauthorized, "failed to verify token")
				return
			}

			ctx := injectSession(r.Context(), sess)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

func checkAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := extractSession(r.Context())
		if err != nil {
			log.Errorf("failed to extract session from context: %v", err)
			responseError(w, http.StatusInternalServerError, "failed to get session")
			return
		}

		if sess.Privilege != string(user.PrivilegeAdmin) {
			log.Warnf("request without admin privilege")
			responseError(w, http.StatusUnauthorized, "need admin privilege")
			return
		}

		h.ServeHTTP(w, r)
	})
}
