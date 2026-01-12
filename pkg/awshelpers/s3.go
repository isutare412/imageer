package awshelpers

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/isutare412/imageer/pkg/apperr"
)

func WrapS3Error(err error, msg string, args ...any) error {
	if err == nil {
		return nil
	}

	switch {
	case isErrorS3NotFound(err):
		fallthrough
	case isErrorS3NoSuchKey(err):
		return apperr.NewError(apperr.CodeNotFound).
			WithCause(err).
			WithSummary("Resource not found").
			WithDetail(msg, args...)

	default:
		return apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Unexpected S3 error").
			WithDetail(msg, args...)
	}
}

func isErrorS3NoSuchKey(err error) bool {
	var s3Err *types.NoSuchKey
	return errors.As(err, &s3Err)
}

func isErrorS3NotFound(err error) bool {
	var s3Err *types.NotFound
	return errors.As(err, &s3Err)
}
