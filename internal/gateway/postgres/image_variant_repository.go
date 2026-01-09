package postgres

import (
	"context"

	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
	"github.com/isutare412/imageer/pkg/dbhelpers"
)

type ImageVariantRepository struct {
	db *gorm.DB
}

func NewImageVariantRepository(client *Client) *ImageVariantRepository {
	return &ImageVariantRepository{
		db: client.db,
	}
}

func (r *ImageVariantRepository) Create(
	ctx context.Context, variant domain.ImageVariant,
) (domain.ImageVariant, error) {
	tx := GetTxOrDB(ctx, r.db)

	iv := entity.NewImageVariant(variant)
	_, err := gorm.G[entity.Image](tx).
		Where(gen.Image.ID.Eq(iv.ImageID)).
		First(ctx)
	if err != nil {
		return domain.ImageVariant{}, dbhelpers.WrapError(err, "Failed to get image %s", iv.ImageID)
	}

	_, err = gorm.G[entity.Preset](tx).
		Where(gen.Preset.ID.Eq(iv.PresetID)).
		First(ctx)
	if err != nil {
		return domain.ImageVariant{},
			dbhelpers.WrapError(err, "Failed to get preset %s", iv.PresetID)
	}

	if err := gorm.G[entity.ImageVariant](tx).Create(ctx, &iv); err != nil {
		return domain.ImageVariant{}, dbhelpers.WrapError(err, "Failed to create image variant")
	}

	iv, err = gorm.G[entity.ImageVariant](tx).
		Where(gen.ImageVariant.ID.Eq(iv.ID)).
		Preload(gen.ImageVariant.Preset.Name(), nil).
		First(ctx)
	if err != nil {
		return domain.ImageVariant{},
			dbhelpers.WrapError(err, "Failed to fetch created image variant %s", iv.ID)
	}

	return iv.ToDomain(), nil
}
