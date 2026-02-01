package valkey

import (
	"context"

	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/otel/trace"

	"github.com/isutare412/imageer/pkg/dbhelpers"
	"github.com/isutare412/imageer/pkg/tracing"
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
	ctx, span := tracing.StartSpan(ctx, "valkey.ImageNotificationPublisher.PublishUploadDone",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServiceValkey))
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
	ctx, span := tracing.StartSpan(ctx, "valkey.ImageNotificationPublisher.PublishProcessDone",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServiceValkey))
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
