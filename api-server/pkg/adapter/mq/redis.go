package mq

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/isutare412/imageer/api-server/pkg/config"
)

type Redis struct {
	client redis.UniversalClient
}

func (r *Redis) Produce(ctx context.Context, topic string, val []byte) error {
	if _, err := r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: topic,
		ID:     "*", // Use auto generated ID by redis
		Values: []interface{}{"data", val},
	}).Result(); err != nil {
		return fmt.Errorf("on XAdd: %w", err)
	}
	return nil
}

func NewRedis(cfg *config.RedisConfig) (*Redis, error) {
	c := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    cfg.Addrs,
		Password: cfg.Password,
	})

	if err := c.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("on ping redis: %w", err)
	}

	return &Redis{
		client: c,
	}, nil
}
