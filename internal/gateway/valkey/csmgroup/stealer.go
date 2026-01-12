package csmgroup

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/dbhelpers"
	"github.com/isutare412/imageer/pkg/log"
)

type Stealer struct {
	client valkey.Client
	cfg    StealerConfig

	lifetimeCtx    context.Context
	lifetimeCancel context.CancelFunc
	workers        *sync.WaitGroup
}

func NewStealer(client valkey.Client, cfg StealerConfig) *Stealer {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = log.WithAttrContext(ctx)

	return &Stealer{
		client:         client,
		cfg:            cfg,
		lifetimeCtx:    ctx,
		lifetimeCancel: cancel,
		workers:        &sync.WaitGroup{},
	}
}

func (s *Stealer) Run() <-chan Message {
	messageCh := make(chan Message, 1)
	s.workers.Go(func() {
		defer close(messageCh)

		ctx := s.lifetimeCtx
		log.AddArgs(ctx, "consumer", s.cfg.Consumer, "group", s.cfg.Group, "stream", s.cfg.Stream)

		ticker := time.NewTicker(s.cfg.StealInterval)
		defer ticker.Stop()

	STEAL_LOOP:
		for {
			select {
			case <-ticker.C:
			case <-ctx.Done():
				break STEAL_LOOP
			}

			messages, err := s.stealMessages(ctx)
			switch {
			case ctx.Err() != nil:
				break STEAL_LOOP
			case err != nil:
				slog.ErrorContext(ctx, "Failed to steal messages from stream", "error", err)
				continue
			}

			for _, msg := range messages {
				// Send message to channel
				// NOTE: we don't select ctx.Done() here to ensure all messages are sent before shutdown
				messageCh <- msg
			}

			if ctx.Err() != nil {
				break STEAL_LOOP
			}
		}

		slog.InfoContext(ctx, "Valkey message steal loop terminated")
	})

	return messageCh
}

func (s *Stealer) Shutdown() {
	s.lifetimeCancel()
	s.workers.Wait()
}

func (s *Stealer) stealMessages(ctx context.Context) ([]Message, error) {
	// Claim idle messages
	resp := s.client.Do(ctx, s.client.B().Xautoclaim().
		Key(s.cfg.Stream).
		Group(s.cfg.Group).
		Consumer(s.cfg.Consumer).
		MinIdleTime(strconv.Itoa(int(s.cfg.StealMinIdleTime.Milliseconds()))).
		Start("0-0").
		Count(100).
		Build())
	if err := resp.Error(); err != nil {
		return nil, dbhelpers.WrapValkeyError(err, "Failed to autoclaim from stream %s", s.cfg.Stream)
	}

	results, err := resp.ToArray()
	if err != nil {
		return nil, dbhelpers.WrapValkeyError(err, "Failed to parse array")
	}

	entries, err := results[1].AsXRange()
	if err != nil {
		return nil, dbhelpers.WrapValkeyError(err, "Failed to parse xrange format")
	}

	messages := make([]Message, 0, len(entries))
	for _, entry := range entries {
		// Check deliver count and drop if exceeds threshold
		isDropped, err := s.tryDropMessage(ctx, entry)
		switch {
		case err != nil:
			return nil, fmt.Errorf("trying drop message: %w", err)
		case isDropped:
			slog.WarnContext(ctx, "Dropped message due to exceeding max delivery attempt",
				"entryId", entry.ID, "threshold", s.cfg.MaxDeliveryAttempt)
			continue
		}

		slog.InfoContext(ctx, "Stealed message from stream", "entryId", entry.ID)

		msg := []byte(entry.FieldValues["msg"])
		messages = append(messages, Message{
			EntryID: entry.ID,
			Data:    msg,
			Ack: func() error {
				if _, err := s.ackMessage(context.Background(), entry.ID); err != nil {
					return err
				}
				return nil
			},
		})
	}

	return messages, nil
}

func (s *Stealer) tryDropMessage(ctx context.Context, entry valkey.XRangeEntry) (bool, error) {
	deliverCount, err := s.checkDeliverCount(ctx, entry.ID)
	if err != nil {
		return false, fmt.Errorf("checking deliver count: %w", err)
	}

	if deliverCount <= s.cfg.MaxDeliveryAttempt {
		return false, nil
	}

	if _, err := s.ackMessage(ctx, entry.ID); err != nil {
		return false, fmt.Errorf("acking message: %w", err)
	}

	return true, nil
}

func (s *Stealer) checkDeliverCount(ctx context.Context, id string) (int64, error) {
	resp := s.client.Do(ctx, s.client.B().Xpending().
		Key(s.cfg.Stream).
		Group(s.cfg.Group).
		Start(id).
		End(id).
		Count(1).
		Build())
	if err := resp.Error(); err != nil {
		return 0, dbhelpers.WrapValkeyError(err,
			"Failed to XPENDING for entry id %s in stream %s", id, s.cfg.Stream)
	}

	results, err := resp.ToArray()
	switch {
	case err != nil:
		return 0, dbhelpers.WrapValkeyError(err, "Failed to parse array")
	case len(results) == 0:
		return 0, apperr.NewError(apperr.CodeNotFound).
			WithSummary("No pending message found for entry id %s in stream %s", id, s.cfg.Stream)
	}

	pending, err := results[0].ToArray()
	if err != nil {
		return 0, dbhelpers.WrapValkeyError(err, "Failed to parse pending array")
	}

	count, err := pending[3].AsInt64()
	if err != nil {
		return 0, dbhelpers.WrapValkeyError(err, "Failed to parse retry count")
	}

	return count, nil
}

func (s *Stealer) ackMessage(ctx context.Context, id string) (acked bool, err error) {
	resp := s.client.Do(ctx, s.client.B().Xack().
		Key(s.cfg.Stream).
		Group(s.cfg.Group).
		Id(id).
		Build())
	if err := resp.Error(); err != nil {
		return false, dbhelpers.WrapValkeyError(err, "Failed to XACK message %s in stream %s",
			id, s.cfg.Stream)
	}

	acked, err = resp.AsBool()
	if err != nil {
		return false, dbhelpers.WrapValkeyError(err, "Failed to parse bool")
	}

	return acked, nil
}
