package valkey

import (
	"context"

	"github.com/valkey-io/valkey-go"
	"google.golang.org/protobuf/proto"

	"github.com/isutare412/imageer/pkg/apperr"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
)

type ImageEventQueue struct {
	client valkey.Client
	cfg    ImageEventQueueConfig
}

func NewImageEventQueue(cfg ImageEventQueueConfig, c *Client) *ImageEventQueue {
	return &ImageEventQueue{
		client: c.client,
		cfg:    cfg,
	}
}

func (q *ImageEventQueue) PushImageProcessRequest(ctx context.Context,
	req *imageerv1.ImageProcessRequest,
) error {
	reqBytes, err := proto.Marshal(req)
	if err != nil {
		return apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Failed to marshal ImageProcessRequest")
	}

	res := q.client.Do(ctx, q.client.B().Xadd().
		Key(q.cfg.StreamKey).
		Maxlen().Almost().Threshold(q.cfg.StreamSizeString()).
		Id("*").
		FieldValue().
		FieldValue("msg", string(reqBytes)).
		Build())
	if err := res.Error(); err != nil {
		return apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Failed to XADD")
	}

	if _, err := res.ToString(); err != nil {
		return apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Failed to convert response to string")
	}

	return nil
}
