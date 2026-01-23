package valkey

import (
	"context"

	"github.com/valkey-io/valkey-go"
	"google.golang.org/protobuf/proto"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/dbhelpers"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
)

type ImageS3DeleteRequestQueue struct {
	client valkey.Client
	cfg    ImageS3DeleteRequestQueueConfig
}

func NewImageS3DeleteRequestQueue(cfg ImageS3DeleteRequestQueueConfig, c *Client,
) *ImageS3DeleteRequestQueue {
	return &ImageS3DeleteRequestQueue{
		client: c.client,
		cfg:    cfg,
	}
}

func (q *ImageS3DeleteRequestQueue) Push(ctx context.Context, req *imageerv1.ImageS3DeleteRequest,
) error {
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
