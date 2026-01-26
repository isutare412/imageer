package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
	"github.com/isutare412/imageer/pkg/apperr"
)

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			v := recover()
			if v == nil {
				return
			}

			aerr := apperr.NewError(apperr.CodeInternalServerError).
				WithCause(fmt.Errorf("panic recover: %v", v))

			slog.With(
				slog.Any("error", aerr),
				slog.String("stackTrace", aerr.Stack.String()),
			).ErrorContext(r.Context(), "Panic occurred while handling http request")

			gen.RespondError(w, r, aerr)
		}()

		next.ServeHTTP(w, r)
	})
}
