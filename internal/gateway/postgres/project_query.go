package postgres

import (
	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
	"github.com/isutare412/imageer/pkg/dbhelpers"
)

func applyProjectSearchFilter(
	q gorm.ChainInterface[entity.Project], filter domain.ProjectSearchFilter,
) gorm.ChainInterface[entity.Project] {
	if filter.Name != nil {
		q = q.Where(gen.Project.Name.Eq(*filter.Name))
	}
	return q
}

func applyProjectSortFilter(
	q gorm.ChainInterface[entity.Project], filter domain.ProjectSortFilter,
) gorm.ChainInterface[entity.Project] {
	switch {
	case filter.CreatedAt:
		order := gen.Project.CreatedAt.Desc()
		if filter.Direction == dbhelpers.SortDirectionAsc {
			order = gen.Project.CreatedAt.Asc()
		}
		q = q.Order(order)

	case filter.UpdatedAt:
		fallthrough

	default:
		order := gen.Project.UpdatedAt.Desc()
		if filter.Direction == dbhelpers.SortDirectionAsc {
			order = gen.Project.UpdatedAt.Asc()
		}
		q = q.Order(order)
	}

	return q
}
