package middleware

import (
	"net/http"

	"github.com/isutare412/imageer/pkg/tracing"
)

func WithTrace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = tracing.ExtractFromHTTPHeader(ctx, r.Header)

		ctx, span := tracing.StartSpan(ctx, "web.middleware.WithTrace")
		defer span.End()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
