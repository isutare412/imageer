package job

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/isutare412/imageer/image-processor/pkg/config"
	log "github.com/sirupsen/logrus"
)

const (
	topicKey = "prcsQueue"
)

type Service struct {
	mq         MsgQueue
	retryDelay time.Duration
	done       chan struct{}
}

func (s *Service) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(s.retryDelay)
		defer ticker.Stop()

	loop:
		for {
			err := s.consume(ctx)
			if err != nil && !errors.Is(err, context.Canceled) {
				log.Errorf("Failed to consume job: %v", err)
			}

			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				break loop
			}
		}

		s.shutdown()
	}()
}

func (s *Service) shutdown() {
	defer close(s.done)
	log.Infof("Processor service shutdowned successfully")
}

func (s *Service) consume(ctx context.Context) error {
	messages, err := s.mq.Consume(ctx, topicKey, 1)
	if err != nil {
		return fmt.Errorf("on consume mq: %w", err)
	}

	log.Infof("Got messages")
	for msg := range messages {
		log.Infof("Message: %s", string(msg))
	}
	return nil
}

func (s *Service) Done() <-chan struct{} {
	return s.done
}

func NewService(cfg *config.ProcessorConfig, mq MsgQueue) (*Service, error) {
	if err := mq.Init(context.Background(), topicKey); err != nil {
		return nil, fmt.Errorf("on init MQ: %w", err)
	}

	if cfg.RetryDelay < 0 {
		return nil, fmt.Errorf("retry delay should be larger than zero: %d", cfg.RetryDelay)
	}

	return &Service{
		mq:         mq,
		retryDelay: time.Duration(cfg.RetryDelay) * time.Millisecond,
		done:       make(chan struct{}),
	}, nil
}
