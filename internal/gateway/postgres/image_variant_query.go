package postgres

import (
	"gorm.io/gorm/clause"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
)

func buildImageVariantUpdateAssigners(req domain.UpdateImageVariantRequest) []clause.Assigner {
	var assigners []clause.Assigner
	if req.State != nil {
		assigners = append(assigners, gen.ImageVariant.State.Set(*req.State))
	}

	if len(assigners) > 0 {
		assigners = append(assigners, gen.ImageVariant.UpdatedAt.Now())
	}
	return assigners
}
