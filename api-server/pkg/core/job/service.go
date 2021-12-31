package job

import (
	"context"
	"fmt"
)

// TODO: Get from config
const (
	topicKey = "prcsQueue"
)

type Service interface {
	Produce(ctx context.Context, val string) error
}

type service struct {
	mq MsgQueue
}

func (s *service) Produce(ctx context.Context, val string) error {
	if err := s.mq.Produce(ctx, topicKey, []byte(val)); err != nil {
		return fmt.Errorf("on mq produce: %w", err)
	}
	return nil
}

func NewService(mq MsgQueue) Service {
	return &service{
		mq: mq,
	}
}
