package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
	"github.com/isutare412/imageer/pkg/dbhelpers"
)

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(client *Client) *ImageRepository {
	return &ImageRepository{
		db: client.db,
	}
}

func (r *ImageRepository) Create(ctx context.Context, image domain.Image) (domain.Image, error) {
	img := entity.NewImage(image)
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		_, err := gorm.G[entity.Project](tx).
			Where(gen.Project.ID.Eq(img.ProjectID)).
			First(ctx)
		if err != nil {
			return dbhelpers.WrapError(err, "Failed to get project %s", img.ProjectID)
		}

		if err := gorm.G[entity.Image](tx).Create(ctx, &img); err != nil {
			return dbhelpers.WrapError(err, "Failed to create image")
		}

		img, err = gorm.G[entity.Image](tx).
			Where(gen.Image.ID.Eq(img.ID)).
			Preload(gen.Image.Project.Name(), nil).
			Preload(gen.Image.Variants.Name(), nil).
			Preload(gen.Image.Variants.Name()+"."+gen.ImageVariant.Preset.Name(), nil).
			First(ctx)
		if err != nil {
			return dbhelpers.WrapError(err, "Failed to fetch created image %s", img.ID)
		}

		return nil
	}); err != nil {
		return domain.Image{}, fmt.Errorf("during transaction: %w", err)
	}
	return img.ToDomain(), nil
}
