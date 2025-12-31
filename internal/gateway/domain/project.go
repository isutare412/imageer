package domain

import (
	"time"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/pkg/dbhelpers"
)

type Project struct {
	ID              string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Name            string
	Transformations []Transformation
}

type ProjectReference struct {
	ID   string
	Name string
}

type CreateProjectRequest struct {
	Name            string                        `validate:"required,max=128"`
	Transformations []CreateTransformationRequest `validate:"dive,required"`
}

func (r CreateProjectRequest) ToProject() Project {
	return Project{
		Name: r.Name,
		Transformations: lo.Map(r.Transformations, func(t CreateTransformationRequest, _ int) Transformation {
			return t.ToTransformation()
		}),
	}
}

type UpdateProjectRequest struct {
	ID              string                        `validate:"max=36"`
	Name            *string                       `validate:"omitempty,max=128"`
	Transformations []UpsertTransformationRequest `validate:"dive,required"`
}

type Projects struct {
	Items []Project
	Total int64
}

type ListProjectsParams struct {
	Offset *int `validate:"omitempty,min=0"`
	Limit  *int `validate:"omitempty,min=1,max=100"`

	SearchFilter ProjectSearchFilter
	SortFilter   ProjectSortFilter
}

func (p ListProjectsParams) OffsetOrDefault() int {
	return lo.FromPtrOr(p.Offset, 0)
}

func (p ListProjectsParams) LimitOrDefault() int {
	return lo.FromPtrOr(p.Limit, 20)
}

type ProjectSearchFilter struct {
	Name *string
}

type ProjectSortFilter struct {
	CreatedAt bool
	UpdatedAt bool
	Direction dbhelpers.SortDirection
}
