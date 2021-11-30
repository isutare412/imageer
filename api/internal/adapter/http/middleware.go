package http

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type statusLogger struct {
	http.ResponseWriter
	status int
}

func (l *statusLogger) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := statusLogger{ResponseWriter: w, status: http.StatusOK}
		h.ServeHTTP(&logger, r)

		log.Infof("%s - \"%s %s\" %d", r.RemoteAddr, r.Method, r.URL.String(), logger.status)
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
