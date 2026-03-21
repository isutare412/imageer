package kafka

import (
	"context"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/kafkahelpers"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
	"github.com/isutare412/imageer/pkg/tracing"
)

type ImageProcessResultQueue struct {
	client *kgo.Client
	cfg    ImageProcessResultQueueConfig
}

func NewImageProcessResultQueue(cfg ImageProcessResultQueueConfig, client *Client,
) *ImageProcessResultQueue {
	return &ImageProcessResultQueue{
		client: client.inner,
		cfg:    cfg,
	}
}

func (q *ImageProcessResultQueue) Push(ctx context.Context, res *imageerv1.ImageProcessResult,
) error {
	ctx = context.WithoutCancel(ctx)
	ctx, span := tracing.StartSpan(ctx, "kafka.ImageProcessResultQueue.Push",
		trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	if res.TraceContext == nil {
		res.TraceContext = make(map[string]string)
	}
	tracing.InjectToMap(ctx, res.TraceContext)

	data, err := proto.Marshal(res)
	if err != nil {
		return apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Failed to marshal protobuf")
	}

	record := &kgo.Record{
		Topic: q.cfg.Topic,
		Value: data,
	}

	q.client.Produce(ctx, record, func(_ *kgo.Record, err error) {
		if err != nil {
			err = kafkahelpers.WrapKafkaError(err, "Failed to produce image process result")
			slog.Error("Failed to produce image process result",
				"topic", q.cfg.Topic, "error", err)
		}
	})

	return nil
}
