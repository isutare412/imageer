package web

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"

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

		slog.With(
			slog.Int("status", resp.Status),
			slog.String("method", req.Method),
			slog.String("path", req.URL.Path),
			slog.String("remoteAddr", ctx.RealIP()),
			slog.String("userAgent", req.UserAgent()),
			slog.Duration("elapsed", time.Since(before)),
			slog.Int64("reqContentLength", req.ContentLength),
			slog.Int64("respSize", resp.Size),
		).Log(rctx, log.SlogLevelAccess, "Handle HTTP request")

		return err
	}
}
