package valkeypubsub

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/samber/lo"
	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/dbhelpers"
)

var defaultInitialBackoff = 100 * time.Millisecond

// Subscriber holds a dedicated connection for subscribe operations. It is
// important to call [Subscriber.Close] after use to put the connection back to
// the the connection pool.
type Subscriber struct {
	client     valkey.Client
	maxRetries int

	// Theses fields are protected by mu
	mu              *sync.Mutex
	patterns        map[string]struct{}
	channels        map[string]struct{}
	dedicated       valkey.DedicatedClient
	dedicatedWait   chan struct{}
	dedicatedCancel func()

	messageCh      chan valkey.PubSubMessage
	errorCh        chan error
	lifetimeCtx    context.Context
	lifetimeCancel context.CancelFunc
}

func NewSubscriber(client valkey.Client, maxRetries int) *Subscriber {
	ctx, cancel := context.WithCancel(context.Background())

	mu := &sync.Mutex{}

	s := &Subscriber{
		client:         client,
		maxRetries:     maxRetries,
		mu:             mu,
		patterns:       make(map[string]struct{}),
		channels:       make(map[string]struct{}),
		dedicatedWait:  make(chan struct{}),
		messageCh:      make(chan valkey.PubSubMessage, 1),
		errorCh:        make(chan error, 1),
		lifetimeCtx:    ctx,
		lifetimeCancel: cancel,
	}

	go s.maintainConnection()

	return s
}

// Close closes the subscriber and releases the dedicated connection to the
// pool.
func (s *Subscriber) Close() {
	s.lifetimeCancel()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.dedicated != nil {
		s.dedicatedCancel()
	}

	s.dedicated = nil
	s.dedicatedCancel = nil
	s.dedicatedWait = make(chan struct{})
}

// Errors returns a channel that receives fatal errors occurred during
// subscription. You **SHOULD** always read from this channel to avoid
// goroutine leaks. When an error is received, the subscriber is closed
// automatically.
func (s *Subscriber) Errors() <-chan error {
	return s.errorCh
}

// Messages returns a channel that receives messages from subscribed channels.
// You **SHOULD** always read from this channel to avoid goroutine leaks.
//
// Multiple calls to this method return the same channel. So if you want to
// fanout, consider using a separate goroutine to read from this channel and
// send to multiple channels.
func (s *Subscriber) Messages() <-chan valkey.PubSubMessage {
	return s.messageCh
}

// Subscribe subscribes to the given channels. You can receive messages from
// [Subscriber.Messages] channel.
func (s *Subscriber) Subscribe(ctx context.Context, channel ...string) error {
	if len(channel) == 0 {
		return apperr.NewError(apperr.CodeBadRequest).
			WithSummary("channels should not be empty")
	}

	select {
	case <-s.lifetimeCtx.Done():
		return apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("Subscriber is closed")
	default:
	}

	dedicated, err := s.getDedicatedBlocking(ctx)
	if err != nil {
		return fmt.Errorf("getting dedicated: %w", err)
	}

	resp := dedicated.Do(ctx, dedicated.B().Subscribe().
		Channel(channel...).
		Build())
	if err := resp.Error(); err != nil {
		return dbhelpers.WrapValkeyError(err, "Failed to SUBSCRIBE %v", channel)
	}

	s.saveChannels(channel)

	return nil
}

// Psubscribe subscribes to the given channel patterns. You can receive messages
// from [Subscriber.Messages] channel.
func (s *Subscriber) Psubscribe(ctx context.Context, pattern ...string) error {
	if len(pattern) == 0 {
		return apperr.NewError(apperr.CodeBadRequest).
			WithSummary("patterns should not be empty")
	}

	select {
	case <-s.lifetimeCtx.Done():
		return apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("Subscriber is closed")
	default:
	}

	dedicated, err := s.getDedicatedBlocking(ctx)
	if err != nil {
		return fmt.Errorf("getting dedicated: %w", err)
	}

	resp := dedicated.Do(ctx, dedicated.B().Psubscribe().
		Pattern(pattern...).
		Build())
	if err := resp.Error(); err != nil {
		return dbhelpers.WrapValkeyError(err, "Failed to PSUBSCRIBE %v", pattern)
	}

	s.savePatterns(pattern)

	return nil
}

// Unsubscribe unsubscribes from the given channels.
func (s *Subscriber) Unsubscribe(ctx context.Context, channel ...string) error {
	dedicated, err := s.getDedicatedBlocking(ctx)
	if err != nil {
		return fmt.Errorf("getting dedicated: %w", err)
	}

	resp := dedicated.Do(ctx, dedicated.B().Unsubscribe().
		Channel(channel...).
		Build())
	if err := resp.Error(); err != nil {
		return dbhelpers.WrapValkeyError(err, "Failed to UNSUBSCRIBE %v", channel)
	}

	s.removeChannels(channel)

	return nil
}

// Punsubscribe unsubscribes from the given channel patterns.
func (s *Subscriber) Punsubscribe(ctx context.Context, pattern ...string) error {
	dedicated, err := s.getDedicatedBlocking(ctx)
	if err != nil {
		return fmt.Errorf("getting dedicated: %w", err)
	}

	resp := dedicated.Do(ctx, dedicated.B().Punsubscribe().
		Pattern(pattern...).
		Build())
	if err := resp.Error(); err != nil {
		return dbhelpers.WrapValkeyError(err, "Failed to PUNSUBSCRIBE %v", pattern)
	}

	s.removePatterns(pattern)

	return nil
}

func (s *Subscriber) handleMessage(msg valkey.PubSubMessage) {
	s.messageCh <- msg
}

func (s *Subscriber) saveChannels(channels []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, c := range channels {
		s.channels[c] = struct{}{}
	}
}

func (s *Subscriber) savePatterns(patterns []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range patterns {
		s.patterns[p] = struct{}{}
	}
}

func (s *Subscriber) removeChannels(channels []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// NOTE: Valkey pubsub allows unsubscribing from all channels at once by
	// sending no channel.
	if len(s.channels) == 0 {
		s.channels = make(map[string]struct{})
		return
	}

	for _, c := range channels {
		delete(s.channels, c)
	}
}

func (s *Subscriber) removePatterns(patterns []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// NOTE: Valkey pubsub allows unsubscribing from all patterns at once by
	// sending no pattern.
	if len(s.patterns) == 0 {
		s.patterns = make(map[string]struct{})
		return
	}

	for _, p := range patterns {
		delete(s.patterns, p)
	}
}

func (s *Subscriber) getDedicatedBlocking(ctx context.Context) (valkey.DedicatedClient, error) {
	for {
		s.mu.Lock()

		if s.dedicated != nil {
			defer s.mu.Unlock()
			return s.dedicated, nil
		}

		waitCh := s.dedicatedWait
		s.mu.Unlock()

		select {
		case <-waitCh:
			// As dedicated connection is ready, retry to get dedicated client
			continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (s *Subscriber) maintainConnection() {
	retryCount := 0
	backoff := defaultInitialBackoff

	for {
		select {
		case <-s.lifetimeCtx.Done():
			return
		default:
		}

		_, dedicatedCancel, hookErrorCh, err := s.connectAndSubscribe()
		if err != nil {
			retryCount++

			slog.Error("Failed to reconnect to valkey pubsub",
				"retryCount", retryCount,
				"error", err)

			// Give up after max retries
			if retryCount > s.maxRetries {
				s.errorCh <- dbhelpers.WrapValkeyError(err,
					"Failed to maintain pubsub connection even after %d retries", s.maxRetries)
				s.Close()
				return
			}

			// Wait before retrying
			select {
			case <-s.lifetimeCtx.Done():
				return
			case <-time.After(backoff):
			}

			if backoff < 5*time.Second {
				backoff *= 2
			}
			continue
		}

		retryCount = 0
		backoff = defaultInitialBackoff

		select {
		case <-s.lifetimeCtx.Done():
			dedicatedCancel()
			return

		case err := <-hookErrorCh:
			dedicatedCancel()

			s.mu.Lock()
			s.dedicated = nil
			s.dedicatedCancel = nil
			s.dedicatedWait = make(chan struct{})
			s.mu.Unlock()

			slog.Error("Error from valkey pubsub hook; trying to reconnect", "error", err)
			continue
		}
	}
}

func (s *Subscriber) connectAndSubscribe() (dedicated valkey.DedicatedClient,
	cancel func(), hookErrorCh <-chan error, err error,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dedicated, cancel = s.client.Dedicate()

	hookErrorCh = dedicated.SetPubSubHooks(valkey.PubSubHooks{
		OnMessage: s.handleMessage,
	})

	// Gather subscribe targets
	subscribes := make([]valkey.Completed, 0, len(s.channels))
	for channel := range s.channels {
		subscribes = append(subscribes, dedicated.B().Subscribe().
			Channel(channel).
			Build())
	}
	for pattern := range s.patterns {
		subscribes = append(subscribes, dedicated.B().Psubscribe().
			Pattern(pattern).
			Build())
	}

	// Resubscribe to channels and patterns
	if len(subscribes) > 0 {
		resps := dedicated.DoMulti(s.lifetimeCtx, subscribes...)
		for _, resp := range resps {
			if err := resp.Error(); err != nil {
				// NOTE: hookErrorCh is consumed here to avoid goroutine leaks
				cancel()
				go func() {
					<-hookErrorCh
				}()

				return nil, nil, nil, dbhelpers.WrapValkeyError(err,
					"Failed to resubscribe channels and patterns")
			}
		}

		slog.Info("Resubscribe complete after reconnecting to valkey pubsub",
			"channels", lo.Keys(s.channels),
			"patterns", lo.Keys(s.patterns))
	}

	s.dedicated = dedicated
	s.dedicatedCancel = cancel
	close(s.dedicatedWait)

	return dedicated, cancel, hookErrorCh, nil
}
