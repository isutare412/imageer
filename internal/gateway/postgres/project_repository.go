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
	"github.com/isutare412/imageer/pkg/apperr"
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
) (domain.Projects, error) {
	var (
		projects   []entity.Project
		totalCount int64
	)
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// Fetch projects
		q := gorm.G[entity.Project](r.db).Scopes()
		q = applyProjectSearchFilter(q, params.SearchFilter)
		q = applyProjectSortFilter(q, params.SortFilter)
		q = applyPagination(q, params.LimitOrDefault(), params.OffsetOrDefault())
		_projects, err := q.
			Preload(gen.Project.Transformations.Name(), nil).
			Find(ctx)
		if err != nil {
			return dbhelpers.WrapError(err, "Failed to list projects")
		}

		// Fetch total count
		q = gorm.G[entity.Project](tx).Scopes()
		q = applyProjectSearchFilter(q, params.SearchFilter)
		count, err := q.Count(ctx, "COUNT(1)")
		if err != nil {
			return dbhelpers.WrapError(err, "Failed to count projects")
		}

		projects = _projects
		totalCount = count
		return nil
	}); err != nil {
		return domain.Projects{}, fmt.Errorf("during transaction: %w", err)
	}

	return domain.Projects{
		Items: lo.Map(projects, func(p entity.Project, _ int) domain.Project {
			return p.ToDomain()
		}),
		Total: totalCount,
	}, nil
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
		assigners := buildProjectUpdateAssigners(req)
		if len(assigners) > 0 {
			_, err := gorm.G[entity.Project](tx).
				Where(gen.Project.ID.Eq(req.ID)).
				Set(append(assigners, gen.Project.UpdatedAt.Set(time.Now()))...).
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

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	_, err := gorm.G[entity.Project](r.db).
		Where(gen.Project.ID.Eq(id)).
		Delete(ctx)
	if err != nil {
		return dbhelpers.WrapError(err, "Failed to delete project %s", id)
	}
	return nil
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
			transToCreate = append(transToCreate, entity.NewTransformationFromUpsert(projectID, t))
		}
	}

	// Update transformations
	for _, t := range transToUpdate {
		assigners := buildTransformationUpdateAssigners(t)
		if len(assigners) == 0 || t.ID == nil {
			continue
		}

		count, err := gorm.G[entity.Transformation](tx).
			Where(gen.Transformation.ID.Eq(*t.ID)).
			Set(append(assigners, gen.Transformation.UpdatedAt.Set(time.Now()))...).
			Update(ctx)
		if err != nil {
			return dbhelpers.WrapError(err, "Failed to update transformation %s", *t.ID)
		} else if count == 0 {
			return apperr.NewError(apperr.CodeNotFound).
				WithSummary("Transformation %s not found", *t.ID)
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
