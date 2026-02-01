package middleware

import (
	"log/slog"
	"net/http"

	"github.com/isutare412/imageer/pkg/log"
	"github.com/isutare412/imageer/pkg/trace"
)

func WithTrace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = trace.ExtractFromHTTPHeader(ctx, r.Header)

		ctx, span := trace.StartSpan(ctx, "web.middleware.WithTrace")
		defer span.End()

		// NOTE: If sampling decision is "not sampled", trace id will be zero-value.
		spanCtx := span.SpanContext()
		if traceID := spanCtx.TraceID().String(); traceID != "" {
			log.AddAttrs(ctx, slog.String("traceId", traceID))
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
