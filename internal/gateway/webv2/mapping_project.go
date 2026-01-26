package webv2

import (
	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
)

func ProjectToWeb(p domain.Project) gen.Project {
	return gen.Project{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Name:      p.Name,
		Presets: lo.Map(p.Presets,
			func(t domain.Preset, _ int) gen.Preset {
				return PresetToWeb(t)
			}),
		ImageCount: p.ImageCount,
	}
}

func ProjectsToWeb(projs domain.Projects) gen.Projects {
	return gen.Projects{
		Items: lo.Map(projs.Items, func(p domain.Project, _ int) gen.Project {
			return ProjectToWeb(p)
		}),
		Total: projs.Total,
	}
}

func ProjectReferenceToWeb(r domain.ProjectReference) gen.ProjectReference {
	return gen.ProjectReference{
		ID:   r.ID,
		Name: r.Name,
	}
}

func ListProjectsAdminParamsToDomain(params gen.ListProjectsAdminParams) domain.ListProjectsParams {
	var offset *int
	if params.Offset != nil {
		v := int(*params.Offset)
		offset = &v
	}

	var limit *int
	if params.Limit != nil {
		v := int(*params.Limit)
		limit = &v
	}

	return domain.ListProjectsParams{
		Offset: offset,
		Limit:  limit,
	}
}

func CreateProjectAdminRequestToDomain(req gen.CreateProjectAdminRequest,
) domain.CreateProjectRequest {
	return domain.CreateProjectRequest{
		Name: req.Name,
		Presets: lo.Map(req.Presets,
			func(t gen.CreatePresetRequest, _ int) domain.CreatePresetRequest {
				return CreatePresetRequestToDomain(t)
			}),
	}
}

func UpdateProjectAdminRequestToDomain(projectID string, req gen.UpdateProjectAdminRequest,
) domain.UpdateProjectRequest {
	return domain.UpdateProjectRequest{
		ID:   projectID,
		Name: req.Name,
		Presets: lo.Map(req.Presets,
			func(t gen.UpsertPresetRequest, _ int) domain.UpsertPresetRequest {
				return UpsertPresetRequestToDomain(t)
			}),
	}
}
