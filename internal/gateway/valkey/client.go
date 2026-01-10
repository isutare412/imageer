package valkey

import (
	"fmt"

	"github.com/valkey-io/valkey-go"
)

type Client struct {
	client valkey.Client
}

func NewClient(cfg ClientConfig) (*Client, error) {
	opt := valkey.ClientOption{
		// Activate pipelining for context cancel support
		// ref: https://github.com/valkey-io/valkey-go?tab=readme-ov-file#canceling-a-context-before-its-deadline
		AlwaysPipelining: true,
	}
	cfg.applyToOption(&opt)

	client, err := valkey.NewClient(opt)
	if err != nil {
		return nil, fmt.Errorf("creating valkey client: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Shutdown() {
	c.client.Close()
}
