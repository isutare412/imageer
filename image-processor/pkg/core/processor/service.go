package processor

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	topicKey = "prcsQueue"
)

type Service struct {
	mq   MsgQueue
	done chan struct{}
}

func NewService(mq MsgQueue) (*Service, error) {
	if err := mq.Init(context.Background(), topicKey); err != nil {
		return nil, fmt.Errorf("on init MQ: %w", err)
	}

	return &Service{
		mq:   mq,
		done: make(chan struct{}),
	}, nil
}

func (s *Service) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(s.done)
			return
		default:
			messages, err := s.mq.Read(ctx, topicKey, 1)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					log.Errorf("Failed to read MQ: %v", err)
				}
				continue
			}

			log.Infof("Got messages")
			for msg := range messages {
				log.Infof("Message: %s", string(msg))
			}
		}
	}
}

func (s *Service) Done() <-chan struct{} {
	return s.done
}
