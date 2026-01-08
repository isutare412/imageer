package postgres

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
	"github.com/isutare412/imageer/pkg/dbhelpers"
)

func applyPresetSearchFilter(
	q gorm.ChainInterface[entity.Preset], filter domain.PresetSearchFilter,
) gorm.ChainInterface[entity.Preset] {
	if filter.ProjectID != nil {
		q = q.Where(gen.Preset.ProjectID.Eq(*filter.ProjectID))
	}
	if len(filter.Names) > 0 {
		q = q.Where(gen.Preset.Name.In(filter.Names...))
	}
	return q
}

func applyPresetSortFilter(
	q gorm.ChainInterface[entity.Preset], filter domain.PresetSortFilter,
) gorm.ChainInterface[entity.Preset] {
	switch {
	case filter.CreatedAt:
		order := gen.Preset.CreatedAt.Desc()
		if filter.Direction == dbhelpers.SortDirectionAsc {
			order = gen.Preset.CreatedAt.Asc()
		}
		q = q.Order(order)

	case filter.UpdatedAt:
		fallthrough

	default:
		order := gen.Preset.UpdatedAt.Desc()
		if filter.Direction == dbhelpers.SortDirectionAsc {
			order = gen.Preset.UpdatedAt.Asc()
		}
		q = q.Order(order)
	}

	return q
}

func buildPresetUpdateAssigners(req domain.UpsertPresetRequest) []clause.Assigner {
	var assigners []clause.Assigner
	if req.Name != nil {
		assigners = append(assigners, gen.Preset.Name.Set(*req.Name))
	}
	if req.Default != nil {
		assigners = append(assigners, gen.Preset.Default.Set(*req.Default))
	}
	if req.Width != nil {
		assigners = append(assigners, gen.Preset.Width.Set(*req.Width))
	}
	if req.Height != nil {
		assigners = append(assigners, gen.Preset.Height.Set(*req.Height))
	}
	return assigners
}
