package valkey

import (
	"context"
	"fmt"

	"github.com/isutare412/imageer/pkg/dbhelpers/valkeypubsub"
	"github.com/isutare412/imageer/pkg/trace"
)

type ImageProcessDoneSubscriber struct {
	client *Client
	cfg    ImageProcessDoneSubscriberConfig
}

func NewImageProcessDoneSubscriber(cfg ImageProcessDoneSubscriberConfig, c *Client,
) *ImageProcessDoneSubscriber {
	return &ImageProcessDoneSubscriber{
		client: c,
		cfg:    cfg,
	}
}

func (s *ImageProcessDoneSubscriber) Subscribe(ctx context.Context, imageID string,
) (<-chan struct{}, <-chan error) {
	ctx, span := trace.StartSpan(ctx, "valkey.ImageProcessDoneSubscriber.Subscribe")
	defer span.End()

	channel := imageProcessDoneChannel(s.cfg.ChannelPrefix, imageID)

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
