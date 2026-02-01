package postgres

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity"
	"github.com/isutare412/imageer/internal/gateway/postgres/entity/gen"
	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/dbhelpers"
	"github.com/isutare412/imageer/pkg/tracing"
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
	ctx, span := tracing.StartSpan(ctx, "postgres.ProjectRepository.FindByID",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServicePostgres))
	defer span.End()

	tx := GetTxOrDB(ctx, r.db)

	proj, err := gorm.G[entity.Project](tx).
		Where(gen.Project.ID.Eq(id)).
		Preload(gen.Project.Presets.Name(), nil).
		First(ctx)
	if err != nil {
		return domain.Project{}, dbhelpers.WrapGORMError(err, "Failed to get project %s", id)
	}

	// Fetch image count
	imageCount, err := gorm.G[entity.Image](tx).
		Where(gen.Image.ProjectID.Eq(id)).
		Count(ctx, "COUNT(1)")
	if err != nil {
		return domain.Project{}, dbhelpers.WrapGORMError(err,
			"Failed to count images for project %s", id)
	}

	result := proj.ToDomain()
	result.ImageCount = imageCount
	return result, nil
}

func (r *ProjectRepository) List(
	ctx context.Context, params domain.ListProjectsParams,
) (domain.Projects, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.ProjectRepository.List",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServicePostgres))
	defer span.End()

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

	// Fetch image counts in batch
	projectIDs := lo.Map(projects, func(p entity.Project, _ int) string {
		return p.ID
	})
	imageCounts, err := r.getImageCountsByProjectIDs(ctx, tx, projectIDs)
	if err != nil {
		return domain.Projects{}, err
	}

	return domain.Projects{
		Items: lo.Map(projects, func(p entity.Project, _ int) domain.Project {
			proj := p.ToDomain()
			proj.ImageCount = imageCounts[p.ID]
			return proj
		}),
		Total: totalCount,
	}, nil
}

func (r *ProjectRepository) Create(ctx context.Context, req domain.Project) (domain.Project, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.ProjectRepository.Create",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServicePostgres))
	defer span.End()

	tx := GetTxOrDB(ctx, r.db)

	proj := entity.NewProject(req)
	if err := gorm.G[entity.Project](tx).Create(ctx, &proj); err != nil {
		return domain.Project{}, dbhelpers.WrapGORMError(err, "Failed to create project")
	}

	return proj.ToDomain(), nil
}

func (r *ProjectRepository) Update(ctx context.Context, req domain.UpdateProjectRequest,
) (domain.Project, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.ProjectRepository.Update",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServicePostgres))
	defer span.End()

	tx := GetTxOrDB(ctx, r.db)

	if err := r.syncPresets(ctx, tx, req.ID, req.Presets); err != nil {
		return domain.Project{}, fmt.Errorf("syncing presets: %w", err)
	}

	// Update project fields
	assigners := buildProjectUpdateAssigners(req)
	if len(assigners) > 0 {
		_, err := gorm.G[entity.Project](tx).
			Where(gen.Project.ID.Eq(req.ID)).
			Set(assigners...).
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
	ctx, span := tracing.StartSpan(ctx, "postgres.ProjectRepository.Delete",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServicePostgres))
	defer span.End()

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
	projectID string, upsertReqs []domain.UpsertPresetRequest,
) error {
	ctx, span := tracing.StartSpan(ctx, "postgres.ProjectRepository.syncPresets",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServicePostgres))
	defer span.End()

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
			Set(assigners...).
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

func (r *ProjectRepository) getImageCountsByProjectIDs(ctx context.Context, tx *gorm.DB,
	projectIDs []string,
) (map[string]int64, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.ProjectRepository.getImageCountsByProjectIDs",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServicePostgres))
	defer span.End()

	if len(projectIDs) == 0 {
		return make(map[string]int64), nil
	}

	type countResult struct {
		ProjectID string `gorm:"column:project_id"`
		Count     int64  `gorm:"column:count"`
	}

	var results []countResult
	err := gorm.G[entity.Image](tx).
		Select(gen.Image.ProjectID.Column().Name, "COUNT(1) AS count").
		Where(gen.Image.ProjectID.In(projectIDs...)).
		Group(gen.Image.ProjectID.Column().Name).
		Scan(ctx, &results)
	if err != nil {
		return nil, dbhelpers.WrapGORMError(err, "Failed to count images by project")
	}

	counts := make(map[string]int64, len(projectIDs))
	for _, r := range results {
		counts[r.ProjectID] = r.Count
	}
	return counts, nil
}
