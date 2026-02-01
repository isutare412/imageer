package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"

	"github.com/isutare412/imageer/pkg/tracing"
)

func WithTrace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = tracing.ExtractFromHTTPHeader(ctx, r.Header)

		spanOpts := []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindServer),
		}
		spanCtx := trace.SpanContextFromContext(ctx)
		if !spanCtx.IsValid() || !spanCtx.IsRemote() {
			spanOpts = append(spanOpts, trace.WithAttributes(tracing.PeerServiceInternet))
		}

		ctx, span := tracing.StartSpan(ctx, "web.middleware.WithTrace", spanOpts...)
		defer span.End()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
