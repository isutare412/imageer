package processor

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

func (s *Service) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.shutdown()
				return

			default:
				err := s.readMQ(ctx)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						continue
					}
					log.Errorf("Failed to read MQ: %v", err)
					time.Sleep(s.retryDelay)
					continue
				}
			}
		}
	}()
}

func (s *Service) shutdown() {
	defer close(s.done)
	log.Infof("Processor service shutdowned successfully")
}

func (s *Service) readMQ(ctx context.Context) error {
	messages, err := s.mq.Read(ctx, topicKey, 1)
	if err != nil {
		return fmt.Errorf("on read mq: %w", err)
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
