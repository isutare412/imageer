package mq

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/isutare412/imageer/image-processor/pkg/config"
)

type RedisMq struct {
	client redis.UniversalClient
}

func NewRedisMq(cfg *config.RedisConfig) (*RedisMq, error) {
	c := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    cfg.Addrs,
		Password: cfg.Password,
	})

	if err := c.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("on ping redis: %v", err)
	}

	return &RedisMq{
		client: c,
	}, nil
}
