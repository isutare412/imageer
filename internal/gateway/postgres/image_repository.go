package postgres

import (
	"context"

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
	tx := GetTxOrDB(ctx, r.db)

	img := entity.NewImage(image)
	_, err := gorm.G[entity.Project](tx).
		Where(gen.Project.ID.Eq(img.ProjectID)).
		First(ctx)
	if err != nil {
		return domain.Image{}, dbhelpers.WrapError(err, "Failed to get project %s", img.ProjectID)
	}

	if err := gorm.G[entity.Image](tx).Create(ctx, &img); err != nil {
		return domain.Image{}, dbhelpers.WrapError(err, "Failed to create image")
	}

	img, err = gorm.G[entity.Image](tx).
		Where(gen.Image.ID.Eq(img.ID)).
		Preload(gen.Image.Project.Name(), nil).
		Preload(gen.Image.Variants.Name(), nil).
		Preload(gen.Image.Variants.Name()+"."+gen.ImageVariant.Preset.Name(), nil).
		First(ctx)
	if err != nil {
		return domain.Image{}, dbhelpers.WrapError(err, "Failed to fetch created image %s", img.ID)
	}

	return img.ToDomain(), nil
}
