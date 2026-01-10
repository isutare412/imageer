package image

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/images"
	"github.com/isutare412/imageer/pkg/validation"
)

type Service struct {
	s3Presigner   port.S3Presigner
	transactioner port.Transactioner
	imageRepo     port.ImageRepository
	imageVarRepo  port.ImageVariantRepository
	presetRepo    port.PresetRepository

	cfg Config
}

func NewService(cfg Config, s3Presigner port.S3Presigner, transactioner port.Transactioner,
	imageRepo port.ImageRepository, imageVarRepo port.ImageVariantRepository,
	presetRepo port.PresetRepository,
) *Service {
	return &Service{
		s3Presigner:   s3Presigner,
		transactioner: transactioner,
		imageRepo:     imageRepo,
		imageVarRepo:  imageVarRepo,
		presetRepo:    presetRepo,
		cfg:           cfg,
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
