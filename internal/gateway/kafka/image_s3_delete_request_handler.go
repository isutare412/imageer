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

	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/pkg/apperr"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
	"github.com/isutare412/imageer/pkg/tracing"
)

type ImageS3DeleteRequestHandler struct {
	imageSvc port.ImageService
	consumer *Consumer
	cfg      ImageS3DeleteRequestHandlerConfig
}

func NewImageS3DeleteRequestHandler(
	cfg ImageS3DeleteRequestHandlerConfig,
	imageSvc port.ImageService,
) *ImageS3DeleteRequestHandler {
	return &ImageS3DeleteRequestHandler{
		imageSvc: imageSvc,
		cfg:      cfg,
	}
}

func (h *ImageS3DeleteRequestHandler) SetConsumer(c *Consumer)       { h.consumer = c }
func (h *ImageS3DeleteRequestHandler) RetryTopic() string            { return h.cfg.RetryTopic }
func (h *ImageS3DeleteRequestHandler) MaxRetryAttempt() int          { return h.cfg.MaxRetryAttempt }
func (h *ImageS3DeleteRequestHandler) RetryBaseDelay() time.Duration { return h.cfg.RetryBaseDelay }

func (h *ImageS3DeleteRequestHandler) HandleRecord(ctx context.Context, record *kgo.Record) {
	handleCtx, cancel := context.WithTimeout(ctx, h.cfg.HandleTimeout)
	defer cancel()

	err := h.handleRecordData(handleCtx, record.Value)
	switch {
	case apperr.IsErrorStatusCode(err, http.StatusBadRequest):
		slog.WarnContext(handleCtx, "Invalid image S3 delete request data, dropping message",
			"error", err)
	case err != nil:
		slog.ErrorContext(handleCtx, "Failed to handle image S3 delete request", "error", err)
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

func (h *ImageS3DeleteRequestHandler) handleRecordData(ctx context.Context, data []byte) error {
	req := &imageerv1.ImageS3DeleteRequest{}
	if err := proto.Unmarshal(data, req); err != nil {
		return apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Failed to unmarshal image S3 delete request").
			WithCause(err)
	}

	ctx = tracing.ExtractFromMap(ctx, req.TraceContext)
	ctx, span := tracing.StartSpan(ctx, "kafka.ImageS3DeleteRequestHandler.handleRecordData",
		trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	if err := h.imageSvc.DeleteS3Objects(ctx, req); err != nil {
		return fmt.Errorf("deleting S3 objects: %w", err)
	}

	slog.InfoContext(ctx, "Deleted S3 objects for image", "imageId", req.ImageId,
		"projectId", req.ProjectId, "deletedCount", len(req.S3Keys))

	return nil
}
