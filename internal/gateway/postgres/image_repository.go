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

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(client *Client) *ImageRepository {
	return &ImageRepository{
		db: client.db,
	}
}

func (r *ImageRepository) FindByID(ctx context.Context, id string) (domain.Image, error) {
	tx := GetTxOrDB(ctx, r.db)

	img, err := r.get(ctx, tx, id)
	if err != nil {
		return domain.Image{}, fmt.Errorf("getting image: %w", err)
	}

	return img.ToDomain(), nil
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

	img, err = r.get(ctx, tx, img.ID)
	if err != nil {
		return domain.Image{}, fmt.Errorf("getting image: %w", err)
	}

	return img.ToDomain(), nil
}

func (r *ImageRepository) Update(ctx context.Context, req domain.UpdateImageRequest,
) (domain.Image, error) {
	tx := GetTxOrDB(ctx, r.db)

	assigners := buildImageUpdateAssigners(req)
	_, err := gorm.G[entity.Image](tx).
		Where(gen.Image.ID.Eq(req.ID)).
		Set(append(assigners, gen.Image.UpdatedAt.Set(time.Now()))...).
		Update(ctx)
	if err != nil {
		return domain.Image{}, dbhelpers.WrapError(err, "Failed to update image %s", req.ID)
	}

	img, err := r.get(ctx, tx, req.ID)
	if err != nil {
		return domain.Image{}, fmt.Errorf("getting image: %w", err)
	}

	return img.ToDomain(), nil
}

func (r *ImageRepository) get(ctx context.Context, tx *gorm.DB, id string) (entity.Image, error) {
	img, err := gorm.G[entity.Image](tx).
		Where(gen.Image.ID.Eq(id)).
		Preload(gen.Image.Project.Name(), nil).
		Preload(gen.Image.Variants.Name(), nil).
		Preload(gen.Image.Variants.Name()+"."+gen.ImageVariant.Preset.Name(), nil).
		First(ctx)
	if err != nil {
		return entity.Image{}, dbhelpers.WrapError(err, "Failed to fetch created image %s", id)
	}
	return img, nil
}
