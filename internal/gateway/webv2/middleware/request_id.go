package middleware

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/isutare412/imageer/pkg/log"
)

const requestIDHeader = "X-Request-ID"

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get or generate request ID
		id := r.Header.Get(requestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}

		// Add requestId to all log entries
		log.AddAttrs(ctx, slog.String("requestId", id))

		w.Header().Set(requestIDHeader, id)
		next.ServeHTTP(w, r)
	})
}
