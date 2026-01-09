package postgres

import (
	"context"

	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
	"github.com/isutare412/imageer/pkg/dbhelpers"
)

type PresetRepository struct {
	db *gorm.DB
}

func NewPresetRepository(client *Client) *PresetRepository {
	return &PresetRepository{
		db: client.db,
	}
}

func (r *PresetRepository) FindByName(ctx context.Context, projectID, name string,
) (domain.Preset, error) {
	tx := GetTxOrDB(ctx, r.db)

	preset, err := gorm.G[entity.Preset](tx).
		Where(gen.Preset.ProjectID.Eq(projectID)).
		Where(gen.Preset.Name.Eq(name)).
		First(ctx)
	if err != nil {
		return domain.Preset{}, dbhelpers.WrapError(err, "Failed to find preset %s", name)
	}

	return preset.ToDomain(), nil
}

func (r *PresetRepository) List(ctx context.Context, params domain.ListPresetsParams,
) ([]domain.Preset, error) {
	tx := GetTxOrDB(ctx, r.db)

	q := gorm.G[entity.Preset](tx).Scopes()
	q = applyPresetSearchFilter(q, params.SearchFilter)
	q = applyPresetSortFilter(q, params.SortFilter)
	q = applyPagination(q, params.LimitOrDefault(), params.OffsetOrDefault())
	presets, err := q.Find(ctx)
	if err != nil {
		return nil, dbhelpers.WrapError(err, "Failed to list presets")
	}

	return lo.Map(presets, func(p entity.Preset, _ int) domain.Preset {
		return p.ToDomain()
	}), nil
}
