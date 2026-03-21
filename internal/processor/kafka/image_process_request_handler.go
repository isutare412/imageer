package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"

	"github.com/isutare412/imageer/internal/processor/port"
	"github.com/isutare412/imageer/pkg/apperr"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
	"github.com/isutare412/imageer/pkg/tracing"
)

type ImageProcessRequestHandler struct {
	imageSvc port.ImageService
	consumer *Consumer
	cfg      ImageProcessRequestHandlerConfig
}

func NewImageProcessRequestHandler(
	cfg ImageProcessRequestHandlerConfig,
	imageSvc port.ImageService,
) *ImageProcessRequestHandler {
	return &ImageProcessRequestHandler{
		imageSvc: imageSvc,
		cfg:      cfg,
	}
}

func (h *ImageProcessRequestHandler) SetConsumer(c *Consumer)       { h.consumer = c }
func (h *ImageProcessRequestHandler) RetryTopic() string            { return h.cfg.RetryTopic }
func (h *ImageProcessRequestHandler) MaxRetryAttempt() int          { return h.cfg.MaxRetryAttempt }
func (h *ImageProcessRequestHandler) RetryBaseDelay() time.Duration { return h.cfg.RetryBaseDelay }

func (h *ImageProcessRequestHandler) HandleRecord(ctx context.Context, record *kgo.Record) {
	handleCtx, cancel := context.WithTimeout(ctx, h.cfg.HandleTimeout)
	defer cancel()

	err := h.handleRecordData(handleCtx, record.Value)
	switch {
	case apperr.IsErrorStatusCode(err, http.StatusNotFound):
		slog.WarnContext(handleCtx, "Referenced resource not found, dropping message", "error", err)
	case apperr.IsErrorStatusCode(err, http.StatusBadRequest):
		slog.WarnContext(handleCtx, "Invalid image process request data, dropping message", "error", err)
	case err != nil:
		slog.ErrorContext(handleCtx, "Failed to handle image process request", "error", err)
		retryCount := parseRetryCount(record)
		nextRetry := retryCount + 1
		if nextRetry > h.cfg.MaxRetryAttempt {
			slog.ErrorContext(handleCtx, "Max retry attempt reached, dropping message",
				"retryCount", retryCount, "maxRetryAttempt", h.cfg.MaxRetryAttempt)
			return
		}
		h.consumer.scheduleRetry(h, record, nextRetry)
	}
}

func (h *ImageProcessRequestHandler) handleRecordData(ctx context.Context, data []byte) error {
	req := &imageerv1.ImageProcessRequest{}
	if err := proto.Unmarshal(data, req); err != nil {
		return apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Failed to unmarshal image process request").
			WithCause(err)
	}

	ctx = tracing.ExtractFromMap(ctx, req.TraceContext)
	ctx, span := tracing.StartSpan(ctx, "kafka.ImageProcessRequestHandler.handleRecordData",
		trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	slog.InfoContext(ctx, "Received image process request",
		"imageId", req.Image.GetId(), "variantId", req.Variant.GetId())

	if err := h.imageSvc.Process(ctx, req); err != nil {
		return fmt.Errorf("processing image: %w", err)
	}

	return nil
}
