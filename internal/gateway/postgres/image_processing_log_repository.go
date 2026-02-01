package postgres

import (
	"context"

	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/pkg/dbhelpers"
	"github.com/isutare412/imageer/pkg/tracing"
)

type ImageProcessingLogRepository struct {
	db *gorm.DB
}

func NewImageProcessingLogRepository(client *Client) *ImageProcessingLogRepository {
	return &ImageProcessingLogRepository{
		db: client.db,
	}
}

func (r *ImageProcessingLogRepository) Create(ctx context.Context, log domain.ImageProcessingLog,
) (domain.ImageProcessingLog, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.ImageProcessingLogRepository.Create")
	defer span.End()

	tx := GetTxOrDB(ctx, r.db)

	procLog := entity.NewImageProcessingLog(log)
	if err := gorm.G[entity.ImageProcessingLog](tx).Create(ctx, &procLog); err != nil {
		return domain.ImageProcessingLog{}, dbhelpers.WrapGORMError(err, "Failed to create image processing log")
	}

	return procLog.ToDomain(), nil
}
