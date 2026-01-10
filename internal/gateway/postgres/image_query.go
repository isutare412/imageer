package postgres

import (
	"gorm.io/gorm/clause"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
)

func buildImageUpdateAssigners(req domain.UpdateImageRequest) []clause.Assigner {
	var assigners []clause.Assigner
	if req.State != nil {
		assigners = append(assigners, gen.Image.State.Set(*req.State))
	}
	return assigners
}
