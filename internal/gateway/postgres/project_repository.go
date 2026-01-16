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
	tx := GetTxOrDB(ctx, r.db)

	proj, err := gorm.G[entity.Project](tx).
		Where(gen.Project.ID.Eq(id)).
		Preload(gen.Project.Presets.Name(), nil).
		First(ctx)
	if err != nil {
		return domain.Project{}, dbhelpers.WrapGORMError(err, "Failed to get project %s", id)
	}

	return proj.ToDomain(), nil
}

func (r *ProjectRepository) List(
	ctx context.Context, params domain.ListProjectsParams,
) (domain.Projects, error) {
	tx := GetTxOrDB(ctx, r.db)

	// Fetch projects
	q := gorm.G[entity.Project](r.db).Scopes()
	q = applyProjectSearchFilter(q, params.SearchFilter)
	q = applyProjectSortFilter(q, params.SortFilter)
	q = applyPagination(q, params.LimitOrDefault(), params.OffsetOrDefault())
	projects, err := q.
		Preload(gen.Project.Presets.Name(), nil).
		Find(ctx)
	if err != nil {
		return domain.Projects{}, dbhelpers.WrapGORMError(err, "Failed to list projects")
	}

	// Fetch total count
	q = gorm.G[entity.Project](tx).Scopes()
	q = applyProjectSearchFilter(q, params.SearchFilter)
	totalCount, err := q.Count(ctx, "COUNT(1)")
	if err != nil {
		return domain.Projects{}, dbhelpers.WrapGORMError(err, "Failed to count projects")
	}

	return domain.Projects{
		Items: lo.Map(projects, func(p entity.Project, _ int) domain.Project {
			return p.ToDomain()
		}),
		Total: totalCount,
	}, nil
}

func (r *ProjectRepository) Create(ctx context.Context, req domain.Project) (domain.Project, error) {
	tx := GetTxOrDB(ctx, r.db)

	proj := entity.NewProject(req)
	if err := gorm.G[entity.Project](tx).Create(ctx, &proj); err != nil {
		return domain.Project{}, dbhelpers.WrapGORMError(err, "Failed to create project")
	}

	return proj.ToDomain(), nil
}

func (r *ProjectRepository) Update(
	ctx context.Context, req domain.UpdateProjectRequest,
) (domain.Project, error) {
	tx := GetTxOrDB(ctx, r.db)

	if err := r.syncPresets(ctx, tx, req.ID, req.Presets); err != nil {
		return domain.Project{}, fmt.Errorf("syncing presets: %w", err)
	}

	// Update project fields
	assigners := buildProjectUpdateAssigners(req)
	if len(assigners) > 0 {
		_, err := gorm.G[entity.Project](tx).
			Where(gen.Project.ID.Eq(req.ID)).
			Set(append(assigners, gen.Project.UpdatedAt.Set(time.Now()))...).
			Update(ctx)
		if err != nil {
			return domain.Project{}, dbhelpers.WrapGORMError(err, "Failed to update project %s", req.ID)
		}
	}

	// Fetch updated project
	proj, err := gorm.G[entity.Project](tx).
		Where(gen.Project.ID.Eq(req.ID)).
		Preload(gen.Project.Presets.Name(), nil).
		First(ctx)
	if err != nil {
		return domain.Project{},
			dbhelpers.WrapGORMError(err, "Failed to fetch updated project %s", req.ID)
	}

	return proj.ToDomain(), nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	tx := GetTxOrDB(ctx, r.db)

	_, err := gorm.G[entity.Project](tx).
		Where(gen.Project.ID.Eq(id)).
		Delete(ctx)
	if err != nil {
		return dbhelpers.WrapGORMError(err, "Failed to delete project %s", id)
	}
	return nil
}

func (*ProjectRepository) syncPresets(ctx context.Context, tx *gorm.DB,
	projectID string,
	upsertReqs []domain.UpsertPresetRequest,
) error {
	presetsToUpdate := make([]domain.UpsertPresetRequest, 0, len(upsertReqs))
	presetsToCreate := make([]entity.Preset, 0, len(upsertReqs))
	for _, t := range upsertReqs {
		if t.IsUpdateRequest() {
			presetsToUpdate = append(presetsToUpdate, t)
		} else {
			presetsToCreate = append(presetsToCreate, entity.NewPresetFromUpsert(projectID, t))
		}
	}

	// Update presets
	for _, t := range presetsToUpdate {
		assigners := buildPresetUpdateAssigners(t)
		if len(assigners) == 0 || t.ID == nil {
			continue
		}

		count, err := gorm.G[entity.Preset](tx).
			Where(gen.Preset.ID.Eq(*t.ID)).
			Set(append(assigners, gen.Preset.UpdatedAt.Set(time.Now()))...).
			Update(ctx)
		if err != nil {
			return dbhelpers.WrapGORMError(err, "Failed to update preset %s", *t.ID)
		} else if count == 0 {
			return apperr.NewError(apperr.CodeNotFound).
				WithSummary("Preset %s not found", *t.ID)
		}
	}

	// Create presets
	err := gorm.G[entity.Preset](tx).
		CreateInBatches(ctx, &presetsToCreate, 10)
	if err != nil {
		return dbhelpers.WrapGORMError(err, "Failed to create presets")
	}

	presetIDs := make([]string, 0, len(presetsToUpdate)+len(presetsToCreate))
	for _, t := range presetsToUpdate {
		presetIDs = append(presetIDs, *t.ID)
	}
	for _, t := range presetsToCreate {
		presetIDs = append(presetIDs, t.ID)
	}

	// Delete removed presets
	_, err = gorm.G[entity.Preset](tx).
		Where(gen.Preset.ProjectID.Eq(projectID)).
		Where(gen.Preset.ID.NotIn(presetIDs...)).
		Delete(ctx)
	if err != nil {
		return dbhelpers.WrapGORMError(err, "Failed to drop presets not in the upsert list")
	}

	return nil
}
