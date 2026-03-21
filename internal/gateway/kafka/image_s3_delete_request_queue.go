package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/kafkahelpers"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
	"github.com/isutare412/imageer/pkg/tracing"
)

type ImageS3DeleteRequestQueue struct {
	client *kgo.Client
	cfg    ImageS3DeleteRequestQueueConfig
}

func NewImageS3DeleteRequestQueue(cfg ImageS3DeleteRequestQueueConfig, client *Client,
) *ImageS3DeleteRequestQueue {
	return &ImageS3DeleteRequestQueue{
		client: client.inner,
		cfg:    cfg,
	}
}

func (q *ImageS3DeleteRequestQueue) Push(ctx context.Context, req *imageerv1.ImageS3DeleteRequest,
) error {
	ctx, span := tracing.StartSpan(ctx, "kafka.ImageS3DeleteRequestQueue.Push",
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

	results := q.client.ProduceSync(ctx, record)
	if err := results.FirstErr(); err != nil {
		return kafkahelpers.WrapKafkaError(err, "Failed to produce image S3 delete request")
	}

	return nil
}
