package job

import (
	"context"
	"fmt"
)

// TODO: Get from config
const (
	topicKey = "prcsQueue"
)

type Service struct {
	mq MsgQueue
}

func (s *Service) Produce(ctx context.Context, val string) error {
	if err := s.mq.Produce(ctx, topicKey, []byte(val)); err != nil {
		return fmt.Errorf("on mq produce: %w", err)
	}
	return nil
}

func NewService(mq MsgQueue) *Service {
	return &Service{
		mq: mq,
	}
}
