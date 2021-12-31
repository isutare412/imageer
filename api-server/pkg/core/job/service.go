package job

import (
	"context"
	"fmt"

	"github.com/isutare412/imageer/api-server/pkg/config"
)

type Service interface {
	Produce(ctx context.Context, val string) error
}

type service struct {
	reqQueueName string
	resQueueName string

	mq MsgQueue
}

func (s *service) Produce(ctx context.Context, val string) error {
	if err := s.mq.Produce(ctx, s.reqQueueName, []byte(val)); err != nil {
		return fmt.Errorf("on mq produce: %w", err)
	}
	return nil
}

func NewService(cfg *config.JobConfig, mq MsgQueue) Service {
	return &service{
		reqQueueName: cfg.Queue.Request,
		resQueueName: cfg.Queue.Response,
		mq:           mq,
	}
}
