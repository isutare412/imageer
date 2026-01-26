package middleware

import (
	"net/http"

	"github.com/isutare412/imageer/internal/gateway/contextbag"
)

func WithContextBag(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = contextbag.WithBag(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
