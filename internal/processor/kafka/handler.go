package kafka

import (
	"context"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Handler interface {
	HandleRecord(context.Context, *kgo.Record)
	RetryTopic() string
	MaxRetryAttempt() int
	RetryBaseDelay() time.Duration
}
