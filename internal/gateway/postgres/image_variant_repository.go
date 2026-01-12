package postgres

import (
	"context"
	"fmt"
	"time"

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
		return domain.ImageVariant{}, dbhelpers.WrapGORMError(err, "Failed to get image %s", iv.ImageID)
	}

	_, err = gorm.G[entity.Preset](tx).
		Where(gen.Preset.ID.Eq(iv.PresetID)).
		First(ctx)
	if err != nil {
		return domain.ImageVariant{},
			dbhelpers.WrapGORMError(err, "Failed to get preset %s", iv.PresetID)
	}

	if err := gorm.G[entity.ImageVariant](tx).Create(ctx, &iv); err != nil {
		return domain.ImageVariant{}, dbhelpers.WrapGORMError(err, "Failed to create image variant")
	}

	iv, err = r.get(ctx, tx, iv.ID)
	if err != nil {
		return domain.ImageVariant{}, fmt.Errorf("getting image variant: %w", err)
	}

	return iv.ToDomain(), nil
}

func (r *ImageVariantRepository) Update(ctx context.Context, req domain.UpdateImageVariantRequest,
) (domain.ImageVariant, error) {
	tx := GetTxOrDB(ctx, r.db)

	assigners := buildImageVariantUpdateAssigners(req)
	_, err := gorm.G[entity.ImageVariant](tx).
		Where(gen.ImageVariant.ID.Eq(req.ID)).
		Set(append(assigners, gen.ImageVariant.UpdatedAt.Set(time.Now()))...).
		Update(ctx)
	if err != nil {
		return domain.ImageVariant{}, dbhelpers.WrapGORMError(err,
			"Failed to update image variant %s", req.ID)
	}

	iv, err := r.get(ctx, tx, req.ID)
	if err != nil {
		return domain.ImageVariant{}, fmt.Errorf("getting image variant: %w", err)
	}

	return iv.ToDomain(), nil
}

func (r *ImageVariantRepository) get(ctx context.Context, tx *gorm.DB, id string,
) (entity.ImageVariant, error) {
	variant, err := gorm.G[entity.ImageVariant](tx).
		Where(gen.ImageVariant.ID.Eq(id)).
		Preload(gen.ImageVariant.Preset.Name(), nil).
		First(ctx)
	if err != nil {
		return entity.ImageVariant{},
			dbhelpers.WrapGORMError(err, "Failed to get image variant %s", id)
	}

	return variant, nil
}
