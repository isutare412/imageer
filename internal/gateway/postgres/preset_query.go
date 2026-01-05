package postgres

import (
	"gorm.io/gorm/clause"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
)

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
