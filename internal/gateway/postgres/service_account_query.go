package postgres

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
	"github.com/isutare412/imageer/pkg/dbhelpers"
)

func applyServiceAccountSearchFilter(
	q gorm.ChainInterface[entity.ServiceAccount], filter domain.ServiceAccountSearchFilter,
) gorm.ChainInterface[entity.ServiceAccount] {
	if filter.Name != nil {
		q = q.Where(gen.ServiceAccount.Name.Eq(*filter.Name))
	}
	return q
}

func applyServiceAccountSortFilter(
	q gorm.ChainInterface[entity.ServiceAccount], filter domain.ServiceAccountSortFilter,
) gorm.ChainInterface[entity.ServiceAccount] {
	switch {
	case filter.CreatedAt:
		order := gen.ServiceAccount.CreatedAt.Desc()
		if filter.Direction == dbhelpers.SortDirectionAsc {
			order = gen.ServiceAccount.CreatedAt.Asc()
		}
		q = q.Order(order)

	case filter.UpdatedAt:
		fallthrough

	default:
		order := gen.ServiceAccount.UpdatedAt.Desc()
		if filter.Direction == dbhelpers.SortDirectionAsc {
			order = gen.ServiceAccount.UpdatedAt.Asc()
		}
		q = q.Order(order)
	}
	return q
}

func buildServiceAccountUpdateAssigners(req domain.UpdateServiceAccountRequest) []clause.Assigner {
	var assigners []clause.Assigner
	if req.Name != nil {
		assigners = append(assigners, gen.ServiceAccount.Name.Set(*req.Name))
	}
	if req.AccessScope != nil {
		assigners = append(assigners, gen.ServiceAccount.AccessScope.Set(*req.AccessScope))
	}
	if req.ExpireAt != nil {
		assigners = append(assigners, gen.ServiceAccount.ExpireAt.Set(*req.ExpireAt))
	}
	return assigners
}
