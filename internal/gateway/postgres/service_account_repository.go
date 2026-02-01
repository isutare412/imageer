package postgres

import (
	"context"

	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
	"github.com/isutare412/imageer/pkg/dbhelpers"
	"github.com/isutare412/imageer/pkg/tracing"
)

type ServiceAccountRepository struct {
	db *gorm.DB
}

func NewServiceAccountRepository(client *Client) *ServiceAccountRepository {
	return &ServiceAccountRepository{
		db: client.db,
	}
}

func (r *ServiceAccountRepository) FindByID(ctx context.Context, id string,
) (domain.ServiceAccount, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.ServiceAccountRepository.FindByID")
	defer span.End()

	tx := GetTxOrDB(ctx, r.db)

	sa, err := gorm.G[entity.ServiceAccount](tx).
		Where(gen.ServiceAccount.ID.Eq(id)).
		Preload(gen.ServiceAccount.Projects.Name(), nil).
		First(ctx)
	if err != nil {
		return domain.ServiceAccount{},
			dbhelpers.WrapGORMError(err, "Failed to fetch service account %s", id)
	}

	return sa.ToDomain(), nil
}

func (r *ServiceAccountRepository) FindByAPIKeyHash(ctx context.Context, hash string,
) (domain.ServiceAccount, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.ServiceAccountRepository.FindByAPIKeyHash")
	defer span.End()

	tx := GetTxOrDB(ctx, r.db)

	sa, err := gorm.G[entity.ServiceAccount](tx).
		Where(gen.ServiceAccount.APIKeyHash.Eq(hash)).
		Preload(gen.ServiceAccount.Projects.Name(), nil).
		First(ctx)
	if err != nil {
		return domain.ServiceAccount{},
			dbhelpers.WrapGORMError(err, "Failed to fetch service account with API key hash")
	}

	return sa.ToDomain(), nil
}

func (r *ServiceAccountRepository) List(ctx context.Context,
	params domain.ListServiceAccountsParams,
) (domain.ServiceAccounts, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.ServiceAccountRepository.List")
	defer span.End()

	tx := GetTxOrDB(ctx, r.db)

	// Fetch service accounts
	q := gorm.G[entity.ServiceAccount](tx).Scopes()
	q = applyServiceAccountSearchFilter(q, params.SearchFilter)
	q = applyServiceAccountSortFilter(q, params.SortFilter)
	q = applyPagination(q, params.LimitOrDefault(), params.OffsetOrDefault())
	accounts, err := q.
		Preload(gen.ServiceAccount.Projects.Name(), nil).
		Find(ctx)
	if err != nil {
		return domain.ServiceAccounts{}, dbhelpers.WrapGORMError(err, "Failed to list service accounts")
	}

	// Fetch total count
	q = gorm.G[entity.ServiceAccount](tx).Scopes()
	q = applyServiceAccountSearchFilter(q, params.SearchFilter)
	count, err := q.Count(ctx, "COUNT(1)")
	if err != nil {
		return domain.ServiceAccounts{},
			dbhelpers.WrapGORMError(err, "Failed to count service accounts")
	}

	return domain.ServiceAccounts{
		Items: lo.Map(accounts, func(sa entity.ServiceAccount, _ int) domain.ServiceAccount {
			return sa.ToDomain()
		}),
		Total: count,
	}, nil
}

func (r *ServiceAccountRepository) Create(ctx context.Context, req domain.ServiceAccount,
) (domain.ServiceAccount, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.ServiceAccountRepository.Create")
	defer span.End()

	tx := GetTxOrDB(ctx, r.db)

	// Create service account record
	sa := entity.NewServiceAccount(req)
	if err := gorm.G[entity.ServiceAccount](tx).Create(ctx, &sa); err != nil {
		return domain.ServiceAccount{}, dbhelpers.WrapGORMError(err, "Failed to create service account")
	}

	// Associate service account to projects
	saProjects := lo.Map(req.Projects,
		func(p domain.ProjectReference, _ int) entity.ServiceAccountProject {
			return entity.ServiceAccountProject{
				ServiceAccountID: sa.ID,
				ProjectID:        p.ID,
			}
		})
	if len(saProjects) > 0 {
		if err := gorm.G[entity.ServiceAccountProject](tx).
			CreateInBatches(ctx, &saProjects, 10); err != nil {
			return domain.ServiceAccount{},
				dbhelpers.WrapGORMError(err, "Failed to associate service account to projects")
		}
	}

	// Fetch created service account with associations
	sa, err := gorm.G[entity.ServiceAccount](tx).
		Where(gen.ServiceAccount.ID.Eq(sa.ID)).
		Preload(gen.ServiceAccount.Projects.Name(), nil).
		First(ctx)
	if err != nil {
		return domain.ServiceAccount{},
			dbhelpers.WrapGORMError(err, "Failed to fetch service account after creation")
	}

	return sa.ToDomain(), nil
}

func (r *ServiceAccountRepository) Update(ctx context.Context,
	req domain.UpdateServiceAccountRequest,
) (domain.ServiceAccount, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.ServiceAccountRepository.Update")
	defer span.End()

	tx := GetTxOrDB(ctx, r.db)

	// Update service account fields
	assigners := buildServiceAccountUpdateAssigners(req)
	_, err := gorm.G[entity.ServiceAccount](tx).
		Where(gen.ServiceAccount.ID.Eq(req.ID)).
		Set(assigners...).
		Update(ctx)
	if err != nil {
		return domain.ServiceAccount{},
			dbhelpers.WrapGORMError(err, "Failed to update service account %s", req.ID)
	}

	// Delete existing project associations
	_, err = gorm.G[entity.ServiceAccountProject](tx).
		Where(gen.ServiceAccountProject.ServiceAccountID.Eq(req.ID)).
		Delete(ctx)
	if err != nil {
		return domain.ServiceAccount{},
			dbhelpers.WrapGORMError(err,
				"Failed to clear existing project associations of service account %s", req.ID)
	}

	// Associate service account to projects
	saProjects := lo.Map(req.ProjectIDs, func(pid string, _ int) entity.ServiceAccountProject {
		return entity.ServiceAccountProject{
			ServiceAccountID: req.ID,
			ProjectID:        pid,
		}
	})
	if len(saProjects) > 0 {
		if err := gorm.G[entity.ServiceAccountProject](tx).
			CreateInBatches(ctx, &saProjects, 10); err != nil {
			return domain.ServiceAccount{},
				dbhelpers.WrapGORMError(err,
					"Failed to associate service account to projects of service account %s", req.ID)
		}
	}

	// Fetch created service account with associations
	sa, err := gorm.G[entity.ServiceAccount](tx).
		Where(gen.ServiceAccount.ID.Eq(req.ID)).
		Preload(gen.ServiceAccount.Projects.Name(), nil).
		First(ctx)
	if err != nil {
		return domain.ServiceAccount{},
			dbhelpers.WrapGORMError(err, "Failed to fetch service account %s after update", req.ID)
	}

	return sa.ToDomain(), nil
}

func (r *ServiceAccountRepository) Delete(ctx context.Context, id string) error {
	ctx, span := tracing.StartSpan(ctx, "postgres.ServiceAccountRepository.Delete")
	defer span.End()

	tx := GetTxOrDB(ctx, r.db)

	if _, err := gorm.G[entity.ServiceAccount](tx).
		Where(gen.ServiceAccount.ID.Eq(id)).
		Delete(ctx); err != nil {
		return dbhelpers.WrapGORMError(err, "Failed to delete service account %s", id)
	}
	return nil
}
