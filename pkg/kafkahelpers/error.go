package kafkahelpers

import (
	"github.com/isutare412/imageer/pkg/apperr"
)

func WrapKafkaError(err error, msg string, args ...any) error {
	if err == nil {
		return nil
	}

	return apperr.NewError(apperr.CodeInternalServerError).
		WithSummary("Unexpected Kafka error").
		WithDetail(msg, args...).
		WithCause(err)
}
