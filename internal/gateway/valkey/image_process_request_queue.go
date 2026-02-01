package valkey

import (
	"context"

	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/dbhelpers"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
	"github.com/isutare412/imageer/pkg/tracing"
)

type ImageProcessRequestQueue struct {
	client valkey.Client
	cfg    ImageProcessRequestQueueConfig
}

func NewImageProcessRequestQueue(cfg ImageProcessRequestQueueConfig, c *Client) *ImageProcessRequestQueue {
	return &ImageProcessRequestQueue{
		client: c.client,
		cfg:    cfg,
	}
}

func (q *ImageProcessRequestQueue) Push(ctx context.Context, req *imageerv1.ImageProcessRequest,
) error {
	ctx, span := tracing.StartSpan(ctx, "valkey.ImageProcessRequestQueue.Push",
		trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	if req.TraceContext == nil {
		req.TraceContext = make(map[string]string)
	}
	tracing.InjectToMap(ctx, req.TraceContext)

	reqBytes, err := proto.Marshal(req)
	if err != nil {
		return apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Failed to marshal protobuf")
	}

	res := q.client.Do(ctx, q.client.B().Xadd().
		Key(q.cfg.StreamKey).
		Maxlen().Almost().Threshold(q.cfg.StreamSizeString()).
		Id("*").
		FieldValue().
		FieldValue("msg", string(reqBytes)).
		Build())
	if err := res.Error(); err != nil {
		return dbhelpers.WrapValkeyError(err, "Failed to XADD")
	}

	if _, err := res.ToString(); err != nil {
		return dbhelpers.WrapValkeyError(err, "Failed to convert response to string")
	}

	return nil
}
