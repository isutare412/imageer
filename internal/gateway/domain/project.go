package domain

import (
	"time"
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
	Name string `validate:"required,max=128"`
}

type UpdateProjectRequest struct {
	Name            *string `validate:"omitempty,max=128"`
	Transformations struct {
		Add    []CreateTransformationRequest `validate:"dive"`
		Update []UpdateTransformationRequest `validate:"dive"`
		Remove []DeleteTransformationRequest `validate:"dive"`
	}
}

type Projects struct {
	Items []Project
	Total int64
}

type ListProjectsParams struct {
	Offset *int64 `validate:"omitempty,min=0"`
	Limit  *int64 `validate:"omitempty,min=1,max=100"`
}
