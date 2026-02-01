package valkey

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"

	"github.com/isutare412/imageer/internal/processor/port"
	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/dbhelpers/valkeystream"
	"github.com/isutare412/imageer/pkg/log"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
	"github.com/isutare412/imageer/pkg/tracing"
)

type ImageProcessRequestHandler struct {
	client   valkey.Client
	reader   *valkeystream.Reader
	stealer  *valkeystream.Stealer
	imageSvc port.ImageService

	messages chan valkeystream.Message

	cfg          ImageProcessRequestHandlerConfig
	consumerName string

	lifetimeCtx    context.Context
	lifetimeCancel context.CancelFunc
	workers        *sync.WaitGroup
}

func NewImageProcessRequestHandler(cfg ImageProcessRequestHandlerConfig, c *Client,
	imageSvc port.ImageService,
) *ImageProcessRequestHandler {
	consumerName := valkeystream.GenerateConsumerName(cfg.GroupName)

	reader := valkeystream.NewReader(c.client, cfg.ToReaderConfig(consumerName))
	stealer := valkeystream.NewStealer(c.client, cfg.ToStealerConfig(consumerName))

	ctx, cancel := context.WithCancel(context.Background())
	ctx = log.WithAttrContext(ctx)

	return &ImageProcessRequestHandler{
		client:         c.client,
		reader:         reader,
		stealer:        stealer,
		imageSvc:       imageSvc,
		messages:       make(chan valkeystream.Message, 1),
		cfg:            cfg,
		consumerName:   consumerName,
		lifetimeCtx:    ctx,
		lifetimeCancel: cancel,
		workers:        &sync.WaitGroup{},
	}
}

func (h *ImageProcessRequestHandler) Initialize(ctx context.Context) error {
	initializer := valkeystream.NewInitializer(h.client, h.cfg.ToInitializerConfig(h.consumerName))
	if err := initializer.Initialize(ctx); err != nil {
		return fmt.Errorf("initializing consumer group: %w", err)
	}

	reaper := valkeystream.NewReaper(h.client, h.cfg.ToReaperConfig())
	if err := reaper.ReapIdleConsumers(ctx); err != nil {
		return fmt.Errorf("reaping idle consumers: %w", err)
	}

	return nil
}

func (h *ImageProcessRequestHandler) Run() {
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

func (h *ImageProcessRequestHandler) Shutdown() {
	h.stealer.Shutdown()
	h.reader.Shutdown()

	h.lifetimeCancel()
	h.workers.Wait()
}

func (h *ImageProcessRequestHandler) handleMessage(msg valkeystream.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.HandleTimeout)
	defer cancel()

	entry := slog.With("entryId", msg.EntryID)

	err := h.handleMessageData(ctx, msg.Data)
	switch {
	case apperr.IsErrorStatusCode(err, http.StatusNotFound):
		entry.WarnContext(ctx, "Referenced resource not found, dropping message", "error", err)
	case apperr.IsErrorStatusCode(err, http.StatusBadRequest):
		entry.WarnContext(ctx, "Invalid image process request data, dropping message", "error", err)
	case err != nil:
		entry.ErrorContext(ctx, "Failed to handle image process request", "error", err)
		return
	}

	if err := msg.Ack(); err != nil {
		entry.ErrorContext(ctx, "Failed to acknowledge image process request", "error", err)
	}
}

func (h *ImageProcessRequestHandler) handleMessageData(ctx context.Context, data []byte) error {
	req := &imageerv1.ImageProcessRequest{}
	if err := proto.Unmarshal(data, req); err != nil {
		return apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Failed to unmarshal image process request").
			WithCause(err)
	}

	ctx = tracing.ExtractFromMap(ctx, req.TraceContext)
	ctx, span := tracing.StartSpan(ctx, "valkey.ImageProcessRequestHandler.handleMessageData",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(tracing.PeerServiceValkey))
	defer span.End()

	slog.InfoContext(ctx, "Received image process request", "imageId", req.Image.Id,
		"variantId", req.Variant.Id, "presetId", req.Preset.Id)

	if err := h.imageSvc.Process(ctx, req); err != nil {
		return fmt.Errorf("processing image: %w", err)
	}

	return nil
}
