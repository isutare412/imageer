package valkey

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/dbhelpers"
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

func (c *Client) HealthCheck(ctx context.Context) error {
	resp := c.client.Do(ctx, c.client.B().Ping().Build())
	if err := resp.Error(); err != nil {
		return dbhelpers.WrapValkeyError(err, "Failed to PING")
	}

	return nil
}

func (c *Client) ComponentName() string {
	return "ValkeyClient"
}
