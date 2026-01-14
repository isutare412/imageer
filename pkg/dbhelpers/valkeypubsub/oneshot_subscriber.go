package valkeypubsub

import (
	"context"
	"errors"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/dbhelpers"
)

// OneShotSubscriber subscribes to a channel and processes a single message. It
// then unsubscribes automatically.
type OneShotSubscriber struct {
	client valkey.Client
}

func NewOneShotSubscriber(client valkey.Client) *OneShotSubscriber {
	return &OneShotSubscriber{
		client: client,
	}
}

func (s *OneShotSubscriber) Subscribe(ctx context.Context, channel ...string,
) (<-chan valkey.PubSubMessage, <-chan error) {
	if len(channel) == 0 {
		panic("channel must not be empty")
	}

	messageCh := make(chan valkey.PubSubMessage, 1)
	errorCh := make(chan error, 1)
	go func() {
		defer func() {
			close(messageCh)
			close(errorCh)
		}()

		// Create a child context to manage the subscription lifecycle
		subscribeCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		err := s.client.Receive(subscribeCtx,
			s.client.B().Subscribe().
				Channel(channel...).
				Build(),
			func(msg valkey.PubSubMessage) {
				go func() {
					messageCh <- msg

					// Cancel the context to unsubscribe after sending the message
					cancel()
				}()
			})
		switch {
		case ctx.Err() != nil:
			errorCh <- ctx.Err()
			return
		case errors.Is(err, valkey.ErrClosing):
			errorCh <- dbhelpers.WrapValkeyError(err, "Client closed while waiting for subscription %v", channel)
			return
		case err != nil:
			errorCh <- dbhelpers.WrapValkeyError(err, "Failed to subscribe channel %v", channel)
			return
		}
	}()

	return messageCh, errorCh
}
