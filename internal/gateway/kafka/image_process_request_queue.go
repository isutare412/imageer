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

type ImageProcessRequestQueue struct {
	client *kgo.Client
	cfg    ImageProcessRequestQueueConfig
}

func NewImageProcessRequestQueue(cfg ImageProcessRequestQueueConfig, client *Client,
) *ImageProcessRequestQueue {
	return &ImageProcessRequestQueue{
		client: client.inner,
		cfg:    cfg,
	}
}

func (q *ImageProcessRequestQueue) Push(ctx context.Context, req *imageerv1.ImageProcessRequest,
) error {
	ctx, span := tracing.StartSpan(ctx, "kafka.ImageProcessRequestQueue.Push",
		trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	if req.TraceContext == nil {
		req.TraceContext = make(map[string]string)
	}
	tracing.InjectToMap(ctx, req.TraceContext)

	data, err := proto.Marshal(req)
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
			err = kafkahelpers.WrapKafkaError(err, "Failed to produce image process request")
			slog.Error("Failed to produce image process request",
				"topic", q.cfg.Topic, "error", err)
		}
	})

	return nil
}
