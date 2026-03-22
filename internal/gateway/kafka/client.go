package kafka

import (
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
)

type Client struct {
	inner *kgo.Client
}

func NewClient(cfg ClientConfig) (*Client, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Addrs...),
		kgo.SASL(plain.Auth{
			User: cfg.User,
			Pass: cfg.Password,
		}.AsMechanism()),
		kgo.ConsumerGroup(cfg.ConsumerGroup),
		kgo.ConsumeTopics(cfg.ConsumeTopics...),
	}

	if cfg.Partitioner != "" {
		opts = append(opts, cfg.Partitioner.KgoOpt())
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("creating kafka client: %w", err)
	}

	return &Client{inner: client}, nil
}

func (c *Client) Shutdown() {
	c.inner.Close()
}
