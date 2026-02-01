package valkey

import (
	"context"

	"github.com/valkey-io/valkey-go"
	"google.golang.org/protobuf/proto"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/dbhelpers"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
	"github.com/isutare412/imageer/pkg/tracing"
)

type ImageProcessResultQueue struct {
	client valkey.Client
	cfg    ImageProcessResultQueueConfig
}

func NewImageProcessResultQueue(cfg ImageProcessResultQueueConfig, c *Client) *ImageProcessResultQueue {
	return &ImageProcessResultQueue{
		client: c.client,
		cfg:    cfg,
	}
}

func (q *ImageProcessResultQueue) Push(ctx context.Context, req *imageerv1.ImageProcessResult,
) error {
	ctx, span := tracing.StartSpan(ctx, "valkey.ImageProcessResultQueue.Push")
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
