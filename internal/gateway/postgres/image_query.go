package postgres

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
	"github.com/isutare412/imageer/pkg/dbhelpers"
)

func applyImageSearchFilter(
	q gorm.ChainInterface[entity.Image], filter domain.ImageSearchFilter,
) gorm.ChainInterface[entity.Image] {
	if filter.State != nil {
		q = q.Where(gen.Image.State.Eq(*filter.State))
	}
	if filter.UpdatedAtBefore != nil {
		q = q.Where(gen.Image.UpdatedAt.Lt(*filter.UpdatedAtBefore))
	}
	return q
}

func applyImageSortFilter(
	q gorm.ChainInterface[entity.Image], filter domain.ImageSortFilter,
) gorm.ChainInterface[entity.Image] {
	switch {
	case filter.CreatedAt:
		order := gen.Image.CreatedAt.Desc()
		if filter.Direction == dbhelpers.SortDirectionAsc {
			order = gen.Image.CreatedAt.Asc()
		}
		q = q.Order(order)

	case filter.UpdatedAt:
		fallthrough

	default:
		order := gen.Image.UpdatedAt.Desc()
		if filter.Direction == dbhelpers.SortDirectionAsc {
			order = gen.Image.UpdatedAt.Asc()
		}
		q = q.Order(order)
	}

	return q
}

func buildImageUpdateAssigners(req domain.UpdateImageRequest) []clause.Assigner {
	var assigners []clause.Assigner
	if req.State != nil {
		assigners = append(assigners, gen.Image.State.Set(*req.State))
	}
	return assigners
}
