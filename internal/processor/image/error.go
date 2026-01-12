package image

import (
	"strings"

	"github.com/isutare412/imageer/pkg/apperr"
)

func wrapBimgError(err error, msg string, args ...any) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "Unsupported image format"):
		return apperr.NewError(apperr.CodeBadRequest).
			WithCause(err).
			WithSummary("Unsupported image format").
			WithDetail(msg, args...)

	case strings.Contains(errMsg, "Unsupported image output type"):
		return apperr.NewError(apperr.CodeBadRequest).
			WithCause(err).
			WithSummary("Unsupported image output format").
			WithDetail(msg, args...)

	default:
		return apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Image processing error").
			WithDetail(msg, args...)
	}
}
