package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/isutare412/imageer/internal/gateway/contextbag"
	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/log"
)

func WithLogAttrContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rctx := r.Context()
		rctx = log.WithAttrContext(rctx)
		next.ServeHTTP(w, r.WithContext(rctx))
	})
}

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		before := time.Now()
		next.ServeHTTP(w, r)

		var (
			status       int
			responseSize int64
		)
		if metrics, ok := GetResponseMetrics(ctx); ok {
			status = metrics.Status
			responseSize = int64(metrics.ResponseSize)
		}

		entry := slog.With(
			slog.Int("status", status),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("query", r.URL.RawQuery),
			slog.String("remoteAddress", r.RemoteAddr),
			slog.String("userAgent", r.UserAgent()),
			slog.Duration("elapsedTime", time.Since(before)),
			slog.Int64("requestContentLength", r.ContentLength),
			slog.Int64("responseSize", responseSize),
		)

		entry = attachAuthenticationInfo(ctx, entry)
		entry.InfoContext(ctx, "Handle HTTP request")
	})
}

func attachAuthenticationInfo(ctx context.Context, entry *slog.Logger) *slog.Logger {
	bag, ok := contextbag.BagFromContext(ctx)
	if !ok {
		return entry
	}

	switch id := bag.Identity.(type) {
	case domain.UserTokenIdentity:
		entry = entry.With(
			"authenticated", true,
			"authenticationMethod", "userToken",
			slog.Group("user",
				"id", id.Payload.UserID,
				"nickname", id.Payload.Nickname,
				"role", id.Payload.Role,
			),
		)

	case domain.ServiceAccountIdentity:
		entry = entry.With(
			"authenticated", true,
			"authenticationMethod", "serviceAccount",
			slog.Group("serviceAccount",
				"id", id.ServiceAccount.ID,
				"name", id.ServiceAccount.Name,
				"accessScope", id.ServiceAccount.AccessScope,
			),
		)
	}

	return entry
}
