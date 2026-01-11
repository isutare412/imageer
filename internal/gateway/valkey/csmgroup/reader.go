package csmgroup

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/log"
)

type Reader struct {
	client valkey.Client
	cfg    ReaderConfig

	lifetimeCtx    context.Context
	lifetimeCancel context.CancelFunc
	workers        *sync.WaitGroup
}

func NewReader(client valkey.Client, cfg ReaderConfig) *Reader {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = log.WithAttrContext(ctx)

	return &Reader{
		client:         client,
		cfg:            cfg,
		lifetimeCtx:    ctx,
		lifetimeCancel: cancel,
		workers:        &sync.WaitGroup{},
	}
}

func (r *Reader) Run() <-chan Message {
	messageCh := make(chan Message, 1)
	r.workers.Go(func() {
		defer close(messageCh)

		ctx := r.lifetimeCtx
		log.AddArgs(ctx, "consumer", r.cfg.Consumer, "group", r.cfg.Group, "stream", r.cfg.Stream)

	READ_LOOP:
		for {
			messages, err := r.readMessages(ctx)
			switch {
			case ctx.Err() != nil:
				break READ_LOOP
			case err != nil:
				slog.ErrorContext(ctx, "Failed to read messages from stream", "error", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for _, msg := range messages {
				// Send message to channel
				// NOTE: we don't select ctx.Done() here to ensure all messages are sent before shutdown
				messageCh <- msg
			}

			if ctx.Err() != nil {
				break READ_LOOP
			}
		}

		slog.InfoContext(ctx, "Valkey message read loop terminated")
	})

	return messageCh
}

func (r *Reader) Shutdown() {
	r.lifetimeCancel()
	r.workers.Wait()
}

func (r *Reader) readMessages(ctx context.Context) ([]Message, error) {
	resp := r.client.Do(ctx, r.client.B().Xreadgroup().
		Group(r.cfg.Group, r.cfg.Consumer).
		Count(r.cfg.ReadBatchSize).
		Block(r.cfg.ReadBlockTimeout.Milliseconds()).
		Streams().
		Key(r.cfg.Stream).
		Id(">").
		Build())
	err := resp.Error()
	switch {
	case valkey.IsValkeyNil(err):
		return nil, nil // No messages
	case err != nil:
		return nil, apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Failed to read messages from stream %s", r.cfg.Stream)
	}

	results, err := resp.AsXRead()
	if err != nil {
		return nil, apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Failed to parse xreadgroup response")
	}

	entries := results[r.cfg.Stream]
	messages := make([]Message, 0, len(entries))
	for _, entry := range entries {
		msg := []byte(entry.FieldValues[r.cfg.EntryFieldKey])

		messages = append(messages, Message{
			Data: msg,
			Ack: func() error {
				if _, err := r.ackMessage(context.Background(), entry.ID); err != nil {
					return err
				}
				return nil
			},
		})
	}

	return messages, nil
}

func (r *Reader) ackMessage(ctx context.Context, id string) (acked bool, err error) {
	resp := r.client.Do(ctx, r.client.B().Xack().
		Key(r.cfg.Stream).
		Group(r.cfg.Group).
		Id(id).
		Build())
	if err := resp.Error(); err != nil {
		return false, apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Failed to ack message %s in stream %s", id, r.cfg.Stream)
	}

	acked, err = resp.AsBool()
	if err != nil {
		return false, apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Failed to parse bool")
	}

	return acked, nil
}
