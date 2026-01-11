package valkey

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/internal/gateway/valkey/csmgroup"
	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/log"
)

type ImageProcessResultHandler struct {
	client  valkey.Client
	reader  *csmgroup.Reader
	stealer *csmgroup.Stealer

	messages chan csmgroup.Message

	cfg          ImageProcessResultHandlerConfig
	consumerName string

	lifetimeCtx    context.Context
	lifetimeCancel context.CancelFunc
	workers        *sync.WaitGroup
}

func NewImageProcessResultHandler(cfg ImageProcessResultHandlerConfig, c *Client,
) *ImageProcessResultHandler {
	consumerName := csmgroup.GenerateConsumerName(cfg.GroupName)

	reader := csmgroup.NewReader(c.client, cfg.ToReaderConfig(consumerName))
	stealer := csmgroup.NewStealer(c.client, cfg.ToStealerConfig(consumerName))

	ctx, cancel := context.WithCancel(context.Background())
	ctx = log.WithAttrContext(ctx)

	return &ImageProcessResultHandler{
		client:         c.client,
		reader:         reader,
		stealer:        stealer,
		messages:       make(chan csmgroup.Message, 1),
		cfg:            cfg,
		consumerName:   consumerName,
		lifetimeCtx:    ctx,
		lifetimeCancel: cancel,
		workers:        &sync.WaitGroup{},
	}
}

func (h *ImageProcessResultHandler) Initialize(ctx context.Context) error {
	initializer := csmgroup.NewInitializer(h.client, h.cfg.StreamKey, h.cfg.GroupName,
		h.consumerName)
	if err := initializer.Initialize(ctx); err != nil {
		return fmt.Errorf("initializing consumer group: %w", err)
	}

	reaper := csmgroup.NewReaper(h.client, h.cfg.StreamKey, h.cfg.GroupName,
		h.cfg.ReapConsumerIdleTime)
	if err := reaper.ReapIdleConsumers(ctx); err != nil {
		return fmt.Errorf("reaping idle consumers: %w", err)
	}

	return nil
}

func (h *ImageProcessResultHandler) Run() {
	stealMessagCh := h.stealer.Run()
	readMessageCh := h.reader.Run()

	// Steal messages from Valkey PEL
	stealDone := make(chan struct{})
	h.workers.Go(func() {
		for msg := range stealMessagCh {
			h.messages <- msg
		}
		close(stealDone)
	})

	// Read messages from Valkey stream
	readDone := make(chan struct{})
	h.workers.Go(func() {
		for msg := range readMessageCh {
			h.messages <- msg
		}
		close(readDone)
	})

	// Close downstream only after all upstreams are closed
	h.workers.Go(func() {
		<-stealDone
		<-readDone
		close(h.messages)
	})

	// Handle messages with multiple workers
	for range h.cfg.HandleConcurrency {
		h.workers.Go(func() {
			for msg := range h.messages {
				h.handleMessage(msg)
			}
		})
	}
}

func (h *ImageProcessResultHandler) Shutdown() {
	h.stealer.Shutdown()
	h.reader.Shutdown()

	h.lifetimeCancel()
	h.workers.Wait()
}

func (h *ImageProcessResultHandler) handleMessage(msg csmgroup.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.HandleTimeout)
	defer cancel()

	entry := slog.With("entryId", msg.EntryID)

	err := h.handleMessageData(ctx, msg.Data)
	switch {
	case apperr.IsErrorStatusCode(err, http.StatusNotFound):
		entry.WarnContext(ctx, "Referenced resource not found, dropping message", "error", err)
	case apperr.IsErrorStatusCode(err, http.StatusBadRequest):
		entry.WarnContext(ctx, "Invalid image process result data, dropping message", "error", err)
	case err != nil:
		entry.ErrorContext(ctx, "Failed to handle image process result", "error", err)
		return
	}

	if err := msg.Ack(); err != nil {
		entry.ErrorContext(ctx, "Failed to acknowledge image process result", "error", err)
	}
}

func (h *ImageProcessResultHandler) handleMessageData(ctx context.Context, data []byte) error {
	slog.DebugContext(ctx, "Will handle image process result message", "data", string(data))

	return nil
}
