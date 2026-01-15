package valkey

import (
	"context"
	"fmt"

	"github.com/isutare412/imageer/pkg/dbhelpers/valkeypubsub"
)

type ImageProcessDoneSubscriber struct {
	subscriber *valkeypubsub.OneShotSubscriber
	cfg        ImageProcessDoneSubscriberConfig
}

func NewImageProcessDoneSubscriber(cfg ImageProcessDoneSubscriberConfig, c *Client,
) *ImageProcessDoneSubscriber {
	return &ImageProcessDoneSubscriber{
		subscriber: valkeypubsub.NewOneShotSubscriber(c.client),
		cfg:        cfg,
	}
}

func (s *ImageProcessDoneSubscriber) Wait(ctx context.Context, imageID string) error {
	channel := imageProcessDoneChannel(s.cfg.ChannelPrefix, imageID)

	messageCh, errorCh := s.subscriber.Subscribe(ctx, channel)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errorCh:
		return fmt.Errorf("subscribing image process done channel: %w", err)
	case <-messageCh:
		return nil
	}
}
