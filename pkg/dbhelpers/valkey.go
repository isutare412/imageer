package dbhelpers

import (
	"errors"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/apperr"
)

func WrapValkeyError(err error, msg string, args ...any) error {
	if err == nil {
		return nil
	}

	switch {
	case valkey.IsValkeyNil(err):
		return apperr.NewError(apperr.CodeNotFound).
			WithSummary("Resource not found").
			WithDetail(msg, args...).
			WithCause(err)
	case valkey.IsValkeyBusyGroup(err):
		return apperr.NewError(apperr.CodeConflict).
			WithSummary("Resource state conflict").
			WithDetail(msg, args...).
			WithCause(err)
	case valkey.IsParseErr(err):
		return apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("Invalid parsing detected").
			WithDetail(msg, args...).
			WithCause(err)
	case errors.Is(err, valkey.ErrClosing):
		return apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("Client is closing").
			WithDetail(msg, args...).
			WithCause(err)
	default:
		return apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("Unexpected Valkey error").
			WithDetail(msg, args...).
			WithCause(err)
	}
}
