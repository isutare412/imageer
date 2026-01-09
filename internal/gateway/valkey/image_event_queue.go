package valkey

import (
	"context"

	"github.com/valkey-io/valkey-go"
)

type ImageEventQueue struct {
	client valkey.Client
}

func NewImageEventQueue(c *Client) *ImageEventQueue {
	return &ImageEventQueue{
		client: c.client,
	}
}

func (q *ImageEventQueue) PushImageUploadedEvent(ctx context.Context) {
}
