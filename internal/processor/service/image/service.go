package image

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/isutare412/imageer/internal/processor/domain"
	"github.com/isutare412/imageer/internal/processor/port"
	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/images"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
)

type Service struct {
	imageProcessor       port.ImageProcessor
	objectStorage        port.ObjectStorage
	imageProcResultQueue port.ImageProcessResultQueue
}

func NewService(imageProcessor port.ImageProcessor, objectStorage port.ObjectStorage,
	imageProcResultQueue port.ImageProcessResultQueue,
) *Service {
	return &Service{
		imageProcessor:       imageProcessor,
		objectStorage:        objectStorage,
		imageProcResultQueue: imageProcResultQueue,
	}
}

func (s *Service) Process(ctx context.Context, req *imageerv1.ImageProcessRequest) error {
	start := time.Now()

	result := &imageerv1.ImageProcessResult{
		ImageId:        req.Image.Id,
		ImageVariantId: req.Variant.Id,
		PresetId:       req.Preset.Id,
	}

	err := s.processImage(ctx, req)
	if aerr, ok := apperr.AsError(err); ok {
		result.IsSuccess = false
		result.ErrorCode = int32(aerr.Code.ID())
		result.ErrorMessage = err.Error()
	} else if err != nil {
		result.IsSuccess = false
		result.ErrorMessage = err.Error()
	} else {
		result.IsSuccess = true
	}

	result.ProcessingTime = durationpb.New(time.Since(start))

	if err := s.imageProcResultQueue.Push(ctx, result); err != nil {
		return fmt.Errorf("pushing image process result: %w", err)
	}

	slog.InfoContext(ctx, "Send image process result", "imageId", result.ImageId,
		"variantId", result.ImageVariantId, "presetId", result.PresetId,
		"isSuccess", result.IsSuccess, "processingTime", result.ProcessingTime)

	return nil
}

func (s *Service) processImage(ctx context.Context, req *imageerv1.ImageProcessRequest) error {
	imageBytes, err := s.objectStorage.Get(ctx, req.Image.S3Key)
	if err != nil {
		return fmt.Errorf("getting original image: %w", err)
	}

	image := domain.RawImage{
		Data:   imageBytes,
		Format: images.NewFormatFromProto(req.Image.Format),
	}
	preset := domain.NewPreset(req.Preset)

	variant, err := s.imageProcessor.Process(ctx, image, preset)
	if err != nil {
		return fmt.Errorf("processing image: %w", err)
	}

	if err := s.objectStorage.Put(ctx, req.Variant.S3Key, variant.Data); err != nil {
		return fmt.Errorf("putting image variant: %w", err)
	}

	return nil
}
