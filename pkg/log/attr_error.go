package log

import (
	"context"
	"log/slog"

	slogmulti "github.com/samber/slog-multi"

	"github.com/isutare412/imageer/pkg/apperr"
)

// newAttrErrorMiddleware creates a middleware that adds [apperr.Error] replaced
// attributes to all log records if present.
func newAttrErrorMiddleware() slogmulti.Middleware {
	return slogmulti.NewInlineMiddleware(
		func(ctx context.Context, level slog.Level, next func(context.Context, slog.Level) bool) bool {
			return next(ctx, level)
		},

		func(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error {
			var (
				stackTrace    string
				errorCodeID   int
				errorCodeName string
			)
			record.Attrs(func(attr slog.Attr) bool {
				if attr.Value.Kind() != slog.KindAny {
					return true
				}

				err, ok := attr.Value.Any().(error)
				if !ok {
					return true
				}

				aerr, ok := apperr.AsError(err)
				if !ok {
					return true
				}

				stackTrace = aerr.Stack.String()
				errorCodeID = aerr.Code.ID()
				errorCodeName = aerr.Code.Name()
				return false
			})

			if stackTrace != "" {
				record.AddAttrs(slog.Int("errorCodeId", errorCodeID))
				record.AddAttrs(slog.String("errorCodeName", errorCodeName))
				record.AddAttrs(slog.String("stackTrace", stackTrace))
			}

			return next(ctx, record)
		},
		func(attrs []slog.Attr, next func([]slog.Attr) slog.Handler) slog.Handler {
			var (
				stackTrace    string
				errorCodeID   int
				errorCodeName string
			)
			for _, attr := range attrs {
				if attr.Value.Kind() != slog.KindAny {
					continue
				}

				err, ok := attr.Value.Any().(error)
				if !ok {
					continue
				}

				aerr, ok := apperr.AsError(err)
				if !ok {
					continue
				}

				stackTrace = aerr.Stack.String()
				errorCodeID = aerr.Code.ID()
				errorCodeName = aerr.Code.Name()
				break
			}

			if stackTrace != "" {
				attrs = append(attrs, slog.Int("errorCodeId", errorCodeID))
				attrs = append(attrs, slog.String("errorCodeName", errorCodeName))
				attrs = append(attrs, slog.String("stackTrace", stackTrace))
			}

			return next(attrs)
		},
		func(name string, next func(string) slog.Handler) slog.Handler { return next(name) },
	)
}
