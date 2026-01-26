package web

import (
	"context"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/internal/gateway/contextbag"
	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/log"
)

func withLogAttrContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		rctx := ctx.Request().Context()
		rctx = log.WithAttrContext(rctx)
		ctx.SetRequest(ctx.Request().WithContext(rctx))
		return next(ctx)
	}
}

func attachRequestID(ctx echo.Context, id string) {
	log.AddAttrs(ctx.Request().Context(), slog.String("requestId", id))
}

func accessLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		before := time.Now()
		err := next(ctx)

		req := ctx.Request()
		resp := ctx.Response()
		rctx := req.Context()

		entry := slog.With(
			slog.Int("status", resp.Status),
			slog.String("method", req.Method),
			slog.String("path", req.URL.Path),
			slog.String("query", req.URL.RawQuery),
			slog.String("remoteAddress", ctx.RealIP()),
			slog.String("userAgent", req.UserAgent()),
			slog.Duration("elapsedTime", time.Since(before)),
			slog.Int64("requestContentLength", req.ContentLength),
			slog.Int64("responseSize", resp.Size),
		)

		entry = attachAuthenticationInfo(rctx, entry)
		entry.InfoContext(rctx, "Handle HTTP request")

		return err
	}
}

func attachAuthenticationInfo(ctx context.Context, entry *slog.Logger) *slog.Logger {
	bag, ok := contextbag.BagFromContext(ctx)
	if !ok {
		return entry
	}

	switch pp := bag.Passport.(type) {
	case domain.UserTokenPassport:
		entry = entry.With(
			"authenticated", true,
			"authenticationMethod", "userToken",
			slog.Group("user",
				"id", pp.Payload.UserID,
				"nickname", pp.Payload.Nickname,
				"role", pp.Payload.Role,
			),
		)

	case domain.ServiceAccountPassport:
		entry = entry.With(
			"authenticated", true,
			"authenticationMethod", "serviceAccount",
			slog.Group("serviceAccount",
				"id", pp.ServiceAccount.ID,
				"name", pp.ServiceAccount.Name,
				"accessScope", pp.ServiceAccount.AccessScope,
			),
		)
	}

	return entry
}
