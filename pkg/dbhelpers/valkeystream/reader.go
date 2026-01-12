package valkeystream

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/dbhelpers"
	"github.com/isutare412/imageer/pkg/log"
)

// Reader reads messages from valkey stream as a consumer.
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
		log.AddArgs(ctx, "consumer", r.cfg.Consumer.Name, "group", r.cfg.Consumer.Group,
			"stream", r.cfg.Consumer.Stream)

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
		Group(r.cfg.Consumer.Group, r.cfg.Consumer.Name).
		Count(r.cfg.ReadBatchSize).
		Block(r.cfg.ReadBlockTimeout.Milliseconds()).
		Streams().
		Key(r.cfg.Consumer.Stream).
		Id(">").
		Build())
	err := dbhelpers.WrapValkeyError(resp.Error(), "Failed to XREADGROUP")
	switch {
	case apperr.IsErrorStatusCode(err, http.StatusNotFound):
		return nil, nil // No messages
	case err != nil:
		return nil, err
	}

	results, err := resp.AsXRead()
	if err != nil {
		return nil, dbhelpers.WrapValkeyError(err, "Failed to parse xreadgroup response")
	}

	entries := results[r.cfg.Consumer.Stream]
	messages := make([]Message, 0, len(entries))
	for _, entry := range entries {
		msg := []byte(entry.FieldValues[r.cfg.EntryFieldKey])

		messages = append(messages, Message{
			EntryID: entry.ID,
			Data:    msg,
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
		Key(r.cfg.Consumer.Stream).
		Group(r.cfg.Consumer.Group).
		Id(id).
		Build())
	if err := resp.Error(); err != nil {
		return false, dbhelpers.WrapValkeyError(err, "Failed to XACK message %s in stream %s", id, r.cfg.Consumer.Stream)
	}

	acked, err = resp.AsBool()
	if err != nil {
		return false, dbhelpers.WrapValkeyError(err, "Failed to parse XACK response as bool")
	}

	return acked, nil
}
