package sqs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/pkg/apperr"
)

type ImageUploadListener struct {
	client   *sqs.Client
	imageSvc port.ImageService

	workers        *sync.WaitGroup
	lifetimeCtx    context.Context
	lifetimeCancel context.CancelFunc

	cfg ImageUploadListenerConfig
}

func NewImageUploadListener(cfg ImageUploadListenerConfig, imageSvc port.ImageService,
) (*ImageUploadListener, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("loading default aws config: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &ImageUploadListener{
		client:         sqs.NewFromConfig(awsCfg),
		imageSvc:       imageSvc,
		workers:        &sync.WaitGroup{},
		lifetimeCtx:    ctx,
		lifetimeCancel: cancel,
		cfg:            cfg,
	}, nil
}

func (l *ImageUploadListener) Run() {
	l.workers.Go(func() {
	MAIN_LOOP:
		for {
			err := l.handleMessageBatch()
			switch {
			case l.lifetimeCtx.Err() != nil:
				break MAIN_LOOP
			case err != nil:
				slog.Error("Failed to handle image upload event message batch", "error", err)
				time.Sleep(5 * time.Second)
				continue
			}
		}

		slog.Info("Image upload listener loop terminated")
	})
}

func (l *ImageUploadListener) Shutdown() {
	l.lifetimeCancel()
	l.workers.Wait()
}

func (l *ImageUploadListener) handleMessageBatch() error {
	ctx, cancel := context.WithTimeout(context.Background(), l.cfg.BatchHandleTimeout)
	defer cancel()

	// Receive messages from SQS
	// NOTE: We use separate lifetime context to ensure ReceiveMessage is cancelled on shutdown
	output, err := l.client.ReceiveMessage(l.lifetimeCtx, &sqs.ReceiveMessageInput{
		QueueUrl:            &l.cfg.QueueURL,
		MaxNumberOfMessages: l.cfg.BatchCount,
		VisibilityTimeout:   int32(l.cfg.VisibilityTimeout.Seconds()),
		WaitTimeSeconds:     int32(l.cfg.PollingWaitTimeout.Seconds()),
	})
	if err != nil {
		return apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Failed to receive image upload events from SQS")
	}

	// Process messages
	delInput := &sqs.DeleteMessageBatchInput{
		QueueUrl: &l.cfg.QueueURL,
	}
	for _, msg := range output.Messages {
		if err := l.handleEvent(ctx, []byte(*msg.Body)); err != nil {
			slog.ErrorContext(ctx, "Failed to handle image upload event", "error", err)
		} else {
			// Prepare to delete successfully processed message
			delInput.Entries = append(delInput.Entries, types.DeleteMessageBatchRequestEntry{
				Id:            msg.MessageId,
				ReceiptHandle: msg.ReceiptHandle,
			})
		}
	}

	if len(delInput.Entries) == 0 {
		return nil
	}

	// Batch delete processed messages
	delOutput, err := l.client.DeleteMessageBatch(ctx, delInput)
	switch {
	case err != nil:
		return apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Failed to batch delete processed image upload event messages from SQS")

	case len(delOutput.Failed) > 0:
		msgs := make([]string, 0, len(delOutput.Failed))
		for _, failure := range delOutput.Failed {
			msgs = append(msgs, fmt.Sprintf("ID=%s Code=%s Message=%s",
				*failure.Id, *failure.Code, *failure.Message))
		}

		return apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("Some image upload event messages failed to delete from SQS").
			WithDetail("%s", strings.Join(msgs, "\n"))
	}

	return nil
}

func (l *ImageUploadListener) handleEvent(ctx context.Context, msgBytes []byte) error {
	var event events.S3Event
	if err := json.Unmarshal(msgBytes, &event); err != nil {
		return apperr.NewError(apperr.CodeInternalServerError).
			WithCause(err).
			WithSummary("Failed to unmarshal S3 event from SQS message")
	}

	var errs error
	for _, record := range event.Records {
		if err := l.handleRecord(ctx, record); err != nil {
			errs = errors.Join(errs, fmt.Errorf("handling record: %w", err))
		}
	}
	return errs
}

func (l *ImageUploadListener) handleRecord(ctx context.Context, record events.S3EventRecord) error {
	ctx, cancel := context.WithTimeout(ctx, l.cfg.HandleTimeout)
	defer cancel()

	s3Key := record.S3.Object.Key
	err := l.imageSvc.StartImageProcessingOnUpload(ctx, s3Key)
	switch {
	case apperr.IsErrorCode(err, apperr.CodeBadRequest):
		fallthrough
	case apperr.IsErrorCode(err, apperr.CodeNotFound):
		slog.WarnContext(ctx, "Skipping image processing for invalid upload",
			"error", err, "s3Key", s3Key)
		return nil

	case err != nil:
		return fmt.Errorf("starting image processing: %w", err)
	}

	return nil
}
