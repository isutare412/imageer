package image

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/images"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
	"github.com/isutare412/imageer/pkg/validation"
)

type Service struct {
	s3Presigner     port.S3Presigner
	transactioner   port.Transactioner
	imageRepo       port.ImageRepository
	imageVarRepo    port.ImageVariantRepository
	presetRepo      port.PresetRepository
	imageEventQueue port.ImageEventQueue

	cfg Config
}

func NewService(cfg Config, s3Presigner port.S3Presigner, transactioner port.Transactioner,
	imageRepo port.ImageRepository, imageVarRepo port.ImageVariantRepository,
	presetRepo port.PresetRepository, imageEventQueue port.ImageEventQueue,
) *Service {
	return &Service{
		s3Presigner:     s3Presigner,
		transactioner:   transactioner,
		imageRepo:       imageRepo,
		imageVarRepo:    imageVarRepo,
		presetRepo:      presetRepo,
		imageEventQueue: imageEventQueue,
		cfg:             cfg,
	}
}

func (s *Service) CreateUploadURL(ctx context.Context, req domain.CreateUploadURLRequest,
) (domain.UploadURL, error) {
	if err := validation.Validate(req); err != nil {
		return domain.UploadURL{}, fmt.Errorf("validating request: %w", err)
	}

	var image domain.Image
	err := s.transactioner.WithTx(ctx, func(ctx context.Context) error {
		// Check presets exist
		params := domain.ListPresetsParams{
			SearchFilter: domain.PresetSearchFilter{
				ProjectID: &req.ProjectID,
				Names:     req.PresetNames,
			},
		}
		presets, err := s.presetRepo.List(ctx, params)
		if err != nil {
			return fmt.Errorf("listing presets: %w", err)
		}
		diffs := findPresetNameDifference(req.PresetNames, presets)
		if len(diffs) > 0 {
			return apperr.NewError(apperr.CodeNotFound).WithSummary("Presets not found: %v", diffs)
		}

		// Create image record
		imageID := uuid.NewString()
		image = domain.Image{
			ID:       imageID,
			FileName: req.FileName,
			Format:   req.Format,
			State:    images.StateWaitingUpload,
			S3Key:    s.imageS3Key(req.ProjectID, imageID, req.Format),
			URL:      s.imagePublicURL(req.ProjectID, imageID, req.Format),
			Project:  domain.ProjectReference{ID: req.ProjectID},
		}
		image, err = s.imageRepo.Create(ctx, image)
		if err != nil {
			return fmt.Errorf("creating image: %w", err)
		}

		// Create image variant records
		for _, preset := range presets {
			variant := domain.ImageVariant{
				Format:  preset.Format,
				State:   images.VariantStateWaitingUpload,
				S3Key:   s.imageVariantS3Key(req.ProjectID, imageID, preset.ID, preset.Format),
				URL:     s.imageVariantPublicURL(req.ProjectID, imageID, preset.ID, preset.Format),
				ImageID: imageID,
				Preset:  domain.PresetReference{ID: preset.ID},
			}
			if _, err = s.imageVarRepo.Create(ctx, variant); err != nil {
				return fmt.Errorf("creating image variant for preset: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return domain.UploadURL{}, fmt.Errorf("during transaction: %w", err)
	}

	// Presign image upload URL
	presignReq := domain.PresignPutObjectRequest{
		S3Key:       image.S3Key,
		ContentType: req.Format.ContentType(),
	}
	presignResp, err := s.s3Presigner.PresignPutObject(ctx, presignReq)
	if err != nil {
		return domain.UploadURL{}, fmt.Errorf("presigning put object: %w", err)
	}

	return domain.UploadURL{
		ImageID:   image.ID,
		ExpiresAt: presignResp.ExpireAt,
		URL:       presignResp.URL,
		Header:    presignResp.Header,
	}, nil
}

func (s *Service) StartImageProcessingOnUpload(ctx context.Context, s3Key string) error {
	_, imageID, ok := parseImageS3Key(s3Key)
	if !ok {
		return apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Unexpected s3 key of uploaded image: %s", s3Key)
	}

	var (
		image        domain.Image
		procRequests []*imageerv1.ImageProcessRequest
	)
	err := s.transactioner.WithTx(ctx, func(ctx context.Context) error {
		var err error

		// Update image state to "ready"
		image, err = s.imageRepo.Update(ctx, domain.UpdateImageRequest{
			ID:    imageID,
			State: lo.ToPtr(images.StateReady),
		})
		if err != nil {
			return fmt.Errorf("updating image: %w", err)
		}

		procRequests = make([]*imageerv1.ImageProcessRequest, 0, len(image.Variants))
		for _, variant := range image.Variants {
			// Update image variant state to "processing"
			variant, err := s.imageVarRepo.Update(ctx, domain.UpdateImageVariantRequest{
				ID:    variant.ID,
				State: lo.ToPtr(images.VariantStateProcessing),
			})
			if err != nil {
				return fmt.Errorf("updating image variant: %w", err)
			}

			preset, err := s.presetRepo.FindByID(ctx, variant.Preset.ID)
			if err != nil {
				return fmt.Errorf("finding preset by ID: %w", err)
			}

			procRequests = append(procRequests, &imageerv1.ImageProcessRequest{
				Image:   image.ToProto(),
				Variant: variant.ToProto(),
				Preset:  preset.ToProto(),
			})
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("during transaction: %w", err)
	}

	for _, req := range procRequests {
		if err := s.imageEventQueue.PushImageProcessRequest(ctx, req); err != nil {
			return fmt.Errorf("enqueuing image process request: %w", err)
		}
	}

	slog.InfoContext(ctx, "Started image processing after client upload", "imageId", image.ID)

	return nil
}
