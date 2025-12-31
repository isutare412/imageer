package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
	"github.com/isutare412/imageer/pkg/dbhelpers"
)

type ServiceAccountRepository struct {
	db *gorm.DB
}

func NewServiceAccountRepository(client *Client) *ServiceAccountRepository {
	return &ServiceAccountRepository{
		db: client.db,
	}
}

func (r *ServiceAccountRepository) FindByID(
	ctx context.Context, id string,
) (domain.ServiceAccount, error) {
	sa, err := gorm.G[entity.ServiceAccount](r.db).
		Where(gen.ServiceAccount.ID.Eq(id)).
		Preload(gen.ServiceAccount.Projects.Name(), nil).
		First(ctx)
	if err != nil {
		return domain.ServiceAccount{},
			dbhelpers.WrapError(err, "Failed to fetch service account %s", id)
	}

	return sa.ToDomain(), nil
}

func (r *ServiceAccountRepository) FindByAPIKeyHash(
	ctx context.Context, hash string,
) (domain.ServiceAccount, error) {
	sa, err := gorm.G[entity.ServiceAccount](r.db).
		Where(gen.ServiceAccount.APIKeyHash.Eq(hash)).
		Preload(gen.ServiceAccount.Projects.Name(), nil).
		First(ctx)
	if err != nil {
		return domain.ServiceAccount{},
			dbhelpers.WrapError(err, "Failed to fetch service account with API key hash")
	}

	return sa.ToDomain(), nil
}

func (r *ServiceAccountRepository) List(
	ctx context.Context, params domain.ListServiceAccountsParams,
) (domain.ServiceAccounts, error) {
	var (
		accounts   []entity.ServiceAccount
		totalCount int64
	)
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// Fetch service accounts
		q := gorm.G[entity.ServiceAccount](tx).Scopes()
		q = applyServiceAccountSearchFilter(q, params.SearchFilter)
		q = applyServiceAccountSortFilter(q, params.SortFilter)
		q = applyPagination(q, params.LimitOrDefault(), params.OffsetOrDefault())
		_accounts, err := q.
			Preload(gen.ServiceAccount.Projects.Name(), nil).
			Find(ctx)
		if err != nil {
			return dbhelpers.WrapError(err, "Failed to list service accounts")
		}

		// Fetch total count
		q = gorm.G[entity.ServiceAccount](tx).Scopes()
		q = applyServiceAccountSearchFilter(q, params.SearchFilter)
		count, err := q.Count(ctx, "COUNT(1)")
		if err != nil {
			return dbhelpers.WrapError(err, "Failed to count service accounts")
		}

		accounts = _accounts
		totalCount = count
		return nil
	}); err != nil {
		return domain.ServiceAccounts{}, fmt.Errorf("during transaction: %w", err)
	}

	return domain.ServiceAccounts{
		Items: lo.Map(accounts, func(sa entity.ServiceAccount, _ int) domain.ServiceAccount {
			return sa.ToDomain()
		}),
		Total: totalCount,
	}, nil
}

func (r *ServiceAccountRepository) Create(
	ctx context.Context, req domain.ServiceAccount,
) (domain.ServiceAccount, error) {
	sa := entity.NewServiceAccount(req)

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// Create service account record
		if err := gorm.G[entity.ServiceAccount](tx).Create(ctx, &sa); err != nil {
			return dbhelpers.WrapError(err, "Failed to create service account")
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
				return dbhelpers.WrapError(err, "Failed to associate service account to projects")
			}
		}

		// Fetch created service account with associations
		saFetched, err := gorm.G[entity.ServiceAccount](tx).
			Where(gen.ServiceAccount.ID.Eq(sa.ID)).
			Preload(gen.ServiceAccount.Projects.Name(), nil).
			First(ctx)
		if err != nil {
			return dbhelpers.WrapError(err, "Failed to fetch service account after creation")
		}

		sa = saFetched
		return nil
	}); err != nil {
		return domain.ServiceAccount{}, fmt.Errorf("during transaction: %w", err)
	}

	return sa.ToDomain(), nil
}

func (r *ServiceAccountRepository) Update(
	ctx context.Context, req domain.UpdateServiceAccountRequest,
) (domain.ServiceAccount, error) {
	var sa entity.ServiceAccount
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// Update service account fields
		assigners := buildServiceAccountUpdateAssigners(req)
		if len(assigners) > 0 {
			_, err := gorm.G[entity.ServiceAccount](tx).
				Where(gen.ServiceAccount.ID.Eq(req.ID)).
				Set(append(assigners, gen.ServiceAccount.UpdatedAt.Set(time.Now()))...).
				Update(ctx)
			if err != nil {
				return dbhelpers.WrapError(err, "Failed to update service account %s", req.ID)
			}
		}

		// Delete existing project associations
		_, err := gorm.G[entity.ServiceAccountProject](tx).
			Where(gen.ServiceAccountProject.ServiceAccountID.Eq(req.ID)).
			Delete(ctx)
		if err != nil {
			return dbhelpers.WrapError(err,
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
				return dbhelpers.WrapError(err,
					"Failed to associate service account to projects of service account %s", req.ID)
			}
		}

		// Fetch created service account with associations
		saFetched, err := gorm.G[entity.ServiceAccount](tx).
			Where(gen.ServiceAccount.ID.Eq(req.ID)).
			Preload(gen.ServiceAccount.Projects.Name(), nil).
			First(ctx)
		if err != nil {
			return dbhelpers.WrapError(err, "Failed to fetch service account %s after update", req.ID)
		}

		sa = saFetched
		return nil
	}); err != nil {
		return domain.ServiceAccount{}, fmt.Errorf("during transaction: %w", err)
	}

	return sa.ToDomain(), nil
}

func (r *ServiceAccountRepository) Delete(ctx context.Context, id string) error {
	if _, err := gorm.G[entity.ServiceAccount](r.db).
		Where(gen.ServiceAccount.ID.Eq(id)).
		Delete(ctx); err != nil {
		return dbhelpers.WrapError(err, "Failed to delete service account %s", id)
	}
	return nil
}
