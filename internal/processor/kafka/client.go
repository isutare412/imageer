package kafka

import (
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Client struct {
	inner *kgo.Client
}

func NewClient(cfg ClientConfig) (*Client, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Addrs...),
		kgo.SASL(cfg.SASLMechanism.Mechanism(cfg.User, cfg.Password)),
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
