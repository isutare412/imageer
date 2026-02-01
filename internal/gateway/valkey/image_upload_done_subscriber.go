package valkey

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"

	"github.com/isutare412/imageer/pkg/dbhelpers/valkeypubsub"
	"github.com/isutare412/imageer/pkg/tracing"
)

type ImageUploadDoneSubscriber struct {
	client *Client
	cfg    ImageUploadDoneSubscriberConfig
}

func NewImageUploadDoneSubscriber(cfg ImageUploadDoneSubscriberConfig, c *Client,
) *ImageUploadDoneSubscriber {
	return &ImageUploadDoneSubscriber{
		client: c,
		cfg:    cfg,
	}
}

func (s *ImageUploadDoneSubscriber) Subscribe(ctx context.Context, imageID string,
) (<-chan struct{}, <-chan error) {
	ctx, span := tracing.StartSpan(ctx, "valkey.ImageUploadDoneSubscriber.Subscribe",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServiceValkey))
	defer span.End()

	channel := imageUploadDoneChannel(s.cfg.ChannelPrefix, imageID)

	notifyCh := make(chan struct{}, 1)
	errorCh := make(chan error, 1)

	sub := valkeypubsub.NewSubscriber(s.client.client, s.cfg.MaxRetries)

	go func() {
		defer close(notifyCh)
		defer close(errorCh)
		defer sub.Close()

		if err := sub.Subscribe(ctx, channel); err != nil {
			errorCh <- fmt.Errorf("subscribing channel: %w", err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case err := <-sub.Errors():
				errorCh <- err
				return
			case <-sub.Messages():
				notifyCh <- struct{}{}
			}
		}
	}()

	return notifyCh, errorCh
}
