package postgres

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
	"github.com/isutare412/imageer/pkg/dbhelpers"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(client *Client) *ProjectRepository {
	return &ProjectRepository{
		db: client.db,
	}
}

func (r *ProjectRepository) FindByID(ctx context.Context, id string) (domain.Project, error) {
	proj, err := gorm.G[entity.Project](r.db).
		Where(gen.Project.ID.Eq(id)).
		Preload(gen.Project.Transformations.Name(), nil).
		First(ctx)
	if err != nil {
		return domain.Project{}, dbhelpers.WrapError(err, "Failed to get project %s", id)
	}

	return proj.ToDomain(), nil
}

func (r *ProjectRepository) List(
	ctx context.Context, params domain.ListProjectsParams,
) ([]domain.Project, error) {
	q := gorm.G[entity.Project](r.db).Scopes()

	// Where clauses
	if params.SearchFilter.Name != nil {
		q = q.Where(gen.Project.Name.Eq(*params.SearchFilter.Name))
	}

	// Order clauses
	switch {
	case params.SortFilter.CreatedAt:
		order := gen.Project.CreatedAt.Asc()
		if params.SortFilter.Direction != dbhelpers.SortDirectionAsc {
			order = gen.Project.CreatedAt.Desc()
		}
		q = q.Order(order)
	case params.SortFilter.UpdatedAt:
		fallthrough
	default:

		order := gen.Project.UpdatedAt.Asc()
		if params.SortFilter.Direction != dbhelpers.SortDirectionAsc {
			order = gen.Project.UpdatedAt.Desc()
		}
		q = q.Order(order)
	}

	// Pagination clauses
	q = q.Offset(params.OffsetOrDefault())
	q = q.Limit(params.LimitOrDefault())

	projects, err := q.
		Preload(gen.Project.Transformations.Name(), nil).
		Find(ctx)
	if err != nil {
		return nil, dbhelpers.WrapError(err, "Failed to list projects")
	}

	return lo.Map(projects, func(p entity.Project, _ int) domain.Project {
		return p.ToDomain()
	}), nil
}

func (r *ProjectRepository) Create(ctx context.Context, req domain.Project) (domain.Project, error) {
	proj := entity.NewProject(req)

	if err := gorm.G[entity.Project](r.db).Create(ctx, &proj); err != nil {
		return domain.Project{}, dbhelpers.WrapError(err, "Failed to create project")
	}

	return proj.ToDomain(), nil
}

func (r *ProjectRepository) Update(
	ctx context.Context, req domain.UpdateProjectRequest,
) (domain.Project, error) {
	var proj domain.Project
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.syncTransformations(ctx, tx, req.ID, req.Transformations); err != nil {
			return fmt.Errorf("syncing transformations: %w", err)
		}

		// Update project fields
		var assigners []clause.Assigner
		if req.Name != nil {
			assigners = append(assigners, gen.Project.Name.Set(*req.Name))
		}
		if len(assigners) > 0 {
			_, err := gorm.G[entity.Project](tx).
				Where(gen.Project.ID.Eq(req.ID)).
				Set(assigners...).
				Update(ctx)
			if err != nil {
				return dbhelpers.WrapError(err, "Failed to update project %s", req.ID)
			}
		}

		// Fetch updated project
		p, err := gorm.G[entity.Project](tx).
			Where(gen.Project.ID.Eq(req.ID)).
			Preload(gen.Project.Transformations.Name(), nil).
			First(ctx)
		if err != nil {
			return dbhelpers.WrapError(err, "Failed to fetch updated project %s", req.ID)
		}

		proj = p.ToDomain()
		return nil
	}); err != nil {
		return domain.Project{}, fmt.Errorf("during transaction: %w", err)
	}

	return proj, nil
}

func (*ProjectRepository) syncTransformations(ctx context.Context, tx *gorm.DB,
	projectID string,
	upsertReqs []domain.UpsertTransformationRequest,
) error {
	transToUpdate := make([]domain.UpsertTransformationRequest, 0, len(upsertReqs))
	transToCreate := make([]entity.Transformation, 0, len(upsertReqs))
	for _, t := range upsertReqs {
		if t.IsUpdateRequest() {
			transToUpdate = append(transToUpdate, t)
		} else {
			transToCreate = append(transToCreate, entity.NewTransformation(projectID, t))
		}
	}

	// Update transformations
	for _, t := range transToUpdate {
		var assigners []clause.Assigner
		if t.Name != nil {
			assigners = append(assigners, gen.Transformation.Name.Set(*t.Name))
		}
		if t.Default != nil {
			assigners = append(assigners, gen.Transformation.Default.Set(*t.Default))
		}
		if t.Width != nil {
			assigners = append(assigners, gen.Transformation.Width.Set(*t.Width))
		}
		if t.Height != nil {
			assigners = append(assigners, gen.Transformation.Height.Set(*t.Height))
		}
		if len(assigners) == 0 || t.ID == nil {
			continue
		}

		_, err := gorm.G[entity.Transformation](tx).
			Where(gen.Transformation.ID.Eq(*t.ID)).
			Set(assigners...).
			Update(ctx)
		if err != nil {
			return dbhelpers.WrapError(err, "Failed to update transformation %s", *t.ID)
		}
	}

	// Create transformations
	err := gorm.G[entity.Transformation](tx).
		CreateInBatches(ctx, &transToCreate, 10)
	if err != nil {
		return dbhelpers.WrapError(err, "Failed to create transformations")
	}

	transIDs := make([]string, 0, len(transToUpdate)+len(transToCreate))
	for _, t := range transToUpdate {
		transIDs = append(transIDs, *t.ID)
	}
	for _, t := range transToCreate {
		transIDs = append(transIDs, t.ID)
	}

	// Delete removed transformations
	_, err = gorm.G[entity.Transformation](tx).
		Where(gen.Transformation.ID.NotIn(transIDs...)).
		Delete(ctx)
	if err != nil {
		return dbhelpers.WrapError(err, "Failed to drop transformations")
	}

	return nil
}
