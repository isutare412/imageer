package valkey

import (
	"context"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/dbhelpers"
	"github.com/isutare412/imageer/pkg/trace"
)

type ImageNotificationPublisher struct {
	client valkey.Client
	cfg    ImageNotificationPublisherConfig
}

func NewImageNotificationPublisher(cfg ImageNotificationPublisherConfig, c *Client,
) *ImageNotificationPublisher {
	return &ImageNotificationPublisher{
		client: c.client,
		cfg:    cfg,
	}
}

func (p *ImageNotificationPublisher) PublishUploadDone(ctx context.Context, imageID string,
) (receiveCount int64, err error) {
	ctx, span := trace.StartSpan(ctx, "valkey.ImageNotificationPublisher.PublishUploadDone")
	defer span.End()

	channel := imageUploadDoneChannel(p.cfg.UploadDoneChannelPrefix, imageID)

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

func (p *ImageNotificationPublisher) PublishProcessDone(ctx context.Context, imageID string,
) (receiveCount int64, err error) {
	ctx, span := trace.StartSpan(ctx, "valkey.ImageNotificationPublisher.PublishProcessDone")
	defer span.End()

	channel := imageProcessDoneChannel(p.cfg.ProcessDoneChannelPrefix, imageID)

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
