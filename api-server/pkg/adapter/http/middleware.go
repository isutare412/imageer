package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/isutare412/imageer/api-server/pkg/core/auth"
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

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := responseLogger{ResponseWriter: w, status: http.StatusOK}
		h.ServeHTTP(&logger, r)

		log.Infof("%s - \"%s %s\" %d %d",
			r.RemoteAddr, r.Method, r.URL.String(), logger.status, logger.length)
	})
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "*")
		h.ServeHTTP(w, r)
	})
}

func authenticate(authSvc auth.Service) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string
			if rawAuth := r.Header.Get("Authorization"); rawAuth != "" {
				authSplit := strings.SplitN(rawAuth, "Bearer ", 2)
				if len(authSplit) < 2 {
					responseError(w, http.StatusBadRequest, "invalid authorization header")
					return
				}
				token = authSplit[1]
			} else {
				cookie, err := r.Cookie("token")
				if err != nil {
					responseError(w, http.StatusBadRequest, "need token from cookie or authorization header")
					return
				}
				token = cookie.Value
			}

			id, err := authSvc.VerifyToken(auth.Token(token))
			if errors.Is(err, auth.ErrTokenExpired) {
				responseError(w, http.StatusInternalServerError, "token expired")
				return
			} else if err != nil {
				log.Errorf("failed to verify token: %v", err)
				responseError(w, http.StatusInternalServerError, "failed to verify token")
				return
			}

			ctx := auth.ContextWithID(r.Context(), id)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}
