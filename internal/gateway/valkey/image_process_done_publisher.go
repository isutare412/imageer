package valkey

import (
	"context"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/dbhelpers"
)

type ImageProcessDonePublisher struct {
	client valkey.Client
	cfg    ImageProcessDonePublisherConfig
}

func NewImageProcessDonePublisher(cfg ImageProcessDonePublisherConfig, c *Client,
) *ImageProcessDonePublisher {
	return &ImageProcessDonePublisher{
		client: c.client,
		cfg:    cfg,
	}
}

func (p *ImageProcessDonePublisher) Publish(ctx context.Context, imageID string,
) (receiveCount int64, err error) {
	channel := imageProcessDoneChannel(p.cfg.ChannelPrefix, imageID)

	resp := p.client.Do(ctx, p.client.B().Publish().
		Channel(channel).
		Message("").
		Build())
	if err := resp.Error(); err != nil {
		return 0, dbhelpers.WrapValkeyError(err, "Failed to PUBLISH to channel %s", channel)
	}

	receiveCount, err = resp.AsInt64()
	if err != nil {
		return 0, dbhelpers.WrapValkeyError(err, "Failed to convert PUBLISH response to int64")
	}

	return receiveCount, nil
}
