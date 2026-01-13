package web

import (
	"cmp"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/pkg/apperr"
)

func respondError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		rctx := ctx.Request().Context()
		err := next(ctx)
		switch {
		case ctx.Response().Committed:
			return err
		case err == nil:
			return nil
		}

		// Default to 500 Internal Server Error
		var (
			statusCode = http.StatusInternalServerError
			appCode    = apperr.DefaultCode(statusCode)
			msg        = http.StatusText(statusCode)
		)

		var (
			aerr *apperr.Error
			herr *echo.HTTPError
		)
		switch {
		// Override if *apperr.Error
		case errors.As(err, &aerr):
			statusCode = aerr.Code.HTTPStatusCode()
			appCode = aerr.Code
			msg = cmp.Or(aerr.ClientMessage(), msg)

		// Override if *echo.HTTPError
		case errors.As(err, &herr):
			statusCode = herr.Code
			appCode = apperr.DefaultCode(statusCode)
			switch v := herr.Message.(type) {
			case string:
				msg = cmp.Or(v, msg)
			case error:
				msg = cmp.Or(v.Error(), msg)
			}
		}

		errorResp := ErrorResponse{
			CodeID:   int64(appCode.ID()),
			CodeName: appCode.Name(),
			Message:  msg,
		}

		entry := slog.With(
			"statusCode", statusCode,
			"error", err,
			"clientMsg", msg,
		)

		switch {
		case statusCode >= 400 && statusCode < 500:
			entry.WarnContext(rctx, "4xx client error occurred")
		case statusCode >= 500:
			entry.ErrorContext(rctx, "5xx server error occurred")
		}

		if err := ctx.JSON(statusCode, errorResp); err != nil {
			slog.ErrorContext(rctx, "Failed to write JSON error response", "error", err)
		}
		return err
	}
}

func recoverPanic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}

			aerr := apperr.NewError(apperr.CodeInternalServerError).
				WithCause(fmt.Errorf("panic recover: %v", r))

			slog.With(
				slog.Any("error", aerr),
				slog.String("stackTrace", aerr.Stack.String()),
			).ErrorContext(ctx.Request().Context(), "Panic occurred while handling http request")

			// Override the error to be returned
			err = aerr
		}()
		return next(ctx)
	}
}
