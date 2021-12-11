package mq

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/isutare412/imageer/image-processor/pkg/config"
)

type Redis struct {
	client       redis.UniversalClient
	groupName    string
	consumerName string
}

func (r *Redis) Init(ctx context.Context, topic string) error {
	if _, err := r.client.XGroupCreateMkStream(
		ctx, topic, r.groupName, "0",
	).Result(); err != nil && !strings.HasPrefix(err.Error(), "BUSYGROUP") {
		return fmt.Errorf("on XGroupCreateMkStream: %w", err)
	}
	return nil
}

func (r *Redis) Consume(ctx context.Context, topic string, limit int64) (<-chan []byte, error) {
	result, err := r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    r.groupName,
		Consumer: r.consumerName,
		NoAck:    true,
		Block:    0,
		Streams:  []string{topic, ">"},
		Count:    limit,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("on XReadGroup: %w", err)
	}

	messages := make(chan []byte, limit)
	go func() {
		defer close(messages)
		for _, msg := range result[0].Messages {
			bytes := []byte(msg.Values["data"].(string))
			messages <- bytes
		}
	}()
	return messages, nil
}

func NewRedis(cfg *config.RedisConfig) (*Redis, error) {
	c := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    cfg.Addrs,
		Password: cfg.Password,
	})

	if err := c.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("on ping redis: %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("on Hostname: %w", err)
	}

	return &Redis{
		client:       c,
		groupName:    cfg.Stream.GroupName,
		consumerName: hostname,
	}, nil
}
