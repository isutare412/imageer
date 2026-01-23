package image

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/pkg/images"
)

type Closer struct {
	transactioner port.Transactioner
	imageRepo     port.ImageRepository
	imageVarRepo  port.ImageVariantRepository
	cfg           CloserConfig

	ticker *time.Ticker
	stopCh chan struct{}
	doneCh chan struct{}
}

func NewCloser(
	cfg CloserConfig,
	transactioner port.Transactioner,
	imageRepo port.ImageRepository,
	imageVariantRepo port.ImageVariantRepository,
) *Closer {
	return &Closer{
		transactioner: transactioner,
		imageRepo:     imageRepo,
		imageVarRepo:  imageVariantRepo,
		cfg:           cfg,
	}
}

func (c *Closer) OnStartedLeading(ctx context.Context) {
	c.ticker = time.NewTicker(c.cfg.CheckInterval)
	c.stopCh = make(chan struct{})
	c.doneCh = make(chan struct{})

	go c.run(ctx)
}

func (c *Closer) OnStoppedLeading() {
	if c.stopCh != nil {
		close(c.stopCh)
		<-c.doneCh
	}
}

func (c *Closer) run(ctx context.Context) {
	defer close(c.doneCh)
	defer c.ticker.Stop()

	for {
		if err := c.closeExpiredImages(); err != nil {
			slog.ErrorContext(ctx, "Failed to close expired images", "error", err)
		}

		select {
		case <-ctx.Done():
			return
		case <-c.stopCh:
			return
		case <-c.ticker.C:
		}
	}
}

func (c *Closer) closeExpiredImages() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	threshold := time.Now().Add(-(c.cfg.CloseThreshold + time.Minute))

	// Find expired images
	params := domain.ListImagesParams{
		Limit: lo.ToPtr(-1),
		SearchFilter: domain.ImageSearchFilter{
			State:           lo.ToPtr(images.StateUploadPending),
			UpdatedAtBefore: &threshold,
		},
	}
	result, err := c.imageRepo.List(ctx, params)
	if err != nil {
		return fmt.Errorf("listing expired images: %w", err)
	}

	if len(result.Items) == 0 {
		return nil
	}

	for _, img := range result.Items {
		err := c.transactioner.WithTx(ctx, func(txCtx context.Context) error {
			_, err := c.imageRepo.Update(txCtx, domain.UpdateImageRequest{
				ID:    img.ID,
				State: lo.ToPtr(images.StateUploadExpired),
			})
			if err != nil {
				return fmt.Errorf("updating image: %w", err)
			}

			for _, variant := range img.Variants {
				_, err := c.imageVarRepo.Update(txCtx, domain.UpdateImageVariantRequest{
					ID:    variant.ID,
					State: lo.ToPtr(images.VariantStateUploadExpired),
				})
				if err != nil {
					return fmt.Errorf("updating image variant %s: %w", variant.ID, err)
				}
			}

			return nil
		})
		if err != nil {
			slog.ErrorContext(ctx, "Failed to close expired image",
				"imageId", img.ID, "error", err)
			continue
		}

		slog.InfoContext(ctx, "Closed expired image", "imageId", img.ID,
			"variantCount", len(img.Variants))
	}

	return nil
}
