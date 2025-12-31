package web

import (
	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

func ProjectToWeb(p domain.Project) Project {
	return Project{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Name:      p.Name,
		Transformations: lo.Map(p.Transformations,
			func(t domain.Transformation, _ int) Transformation {
				return TransformationToWeb(t)
			}),
	}
}

func ProjectsToWeb(projs domain.Projects) Projects {
	return Projects{
		Items: lo.Map(projs.Items, func(p domain.Project, _ int) Project {
			return ProjectToWeb(p)
		}),
		Total: projs.Total,
	}
}

func ProjectReferenceToWeb(r domain.ProjectReference) ProjectReference {
	return ProjectReference{
		ID:   r.ID,
		Name: r.Name,
	}
}

func ListProjectsAdminParamsToDomain(params ListProjectsAdminParams) domain.ListProjectsParams {
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

func CreateProjectAdminRequestToDomain(req CreateProjectAdminRequest) domain.CreateProjectRequest {
	return domain.CreateProjectRequest{
		Name: req.Name,
		Transformations: lo.Map(req.Transformations,
			func(t CreateTransformationRequest, _ int) domain.CreateTransformationRequest {
				return CreateTransformationRequestToDomain(t)
			}),
	}
}

func UpdateProjectAdminRequestToDomain(projectID string, req UpdateProjectAdminRequest,
) domain.UpdateProjectRequest {
	return domain.UpdateProjectRequest{
		ID:   projectID,
		Name: req.Name,
		Transformations: lo.Map(req.Transformations,
			func(t UpsertTransformationRequest, _ int) domain.UpsertTransformationRequest {
				return UpsertTransformationRequestToDomain(t)
			}),
	}
}
