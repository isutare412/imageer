package gen

import (
	"cmp"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/isutare412/imageer/pkg/apperr"
)

func RespondJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("Failed to json encode response", "error", err)
	}
}

func RespondNoContent(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()

	var (
		statusCode = http.StatusInternalServerError
		appCode    = apperr.DefaultCode(statusCode)
		msg        = http.StatusText(statusCode)
	)

	var (
		aerr              *apperr.Error
		cookieParamErr    *UnescapedCookieParamError
		unmarshalParamErr *UnmarshalingParamError
		requiredParamErr  *RequiredParamError
		invalidParamErr   *InvalidParamFormatError
		tooManyParamErr   *TooManyValuesForParamError
		requiredHeaderErr *RequiredHeaderError
	)
	switch {
	case errors.As(err, &aerr):
		statusCode = aerr.Code.HTTPStatusCode()
		appCode = aerr.Code
		msg = cmp.Or(aerr.ClientMessage(), msg)

	case errors.As(err, &cookieParamErr):
		appCode = apperr.CodeUnprocessableEntity
		statusCode = appCode.HTTPStatusCode()
		msg = cookieParamErr.Error()

	case errors.As(err, &unmarshalParamErr):
		appCode = apperr.CodeUnprocessableEntity
		statusCode = appCode.HTTPStatusCode()
		msg = unmarshalParamErr.Error()

	case errors.As(err, &requiredParamErr):
		appCode = apperr.CodeBadRequest
		statusCode = appCode.HTTPStatusCode()
		msg = requiredParamErr.Error()

	case errors.As(err, &invalidParamErr):
		appCode = apperr.CodeUnprocessableEntity
		statusCode = appCode.HTTPStatusCode()
		msg = invalidParamErr.Error()

	case errors.As(err, &tooManyParamErr):
		appCode = apperr.CodeUnprocessableEntity
		statusCode = appCode.HTTPStatusCode()
		msg = tooManyParamErr.Error()

	case errors.As(err, &requiredHeaderErr):
		appCode = apperr.CodeBadRequest
		statusCode = appCode.HTTPStatusCode()
		msg = requiredHeaderErr.Error()
	}

	entry := slog.With(
		"statusCode", statusCode,
		"error", err,
		"clientMsg", msg,
	)
	switch {
	case statusCode >= 400 && statusCode < 500:
		entry.WarnContext(ctx, "4xx client error occurred")
	case statusCode >= 500:
		entry.ErrorContext(ctx, "5xx server error occurred")
	}

	errorResp := ErrorResponse{
		CodeID:   int64(appCode.ID()),
		CodeName: appCode.Name(),
		Message:  msg,
	}

	RespondJSON(w, statusCode, errorResp)
}
