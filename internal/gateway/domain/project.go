package domain

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID              uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Name            string
	Transformations []Transformation
}

type CreateProjectRequest struct {
	Name string `validate:"required,max=128"`
}

type UpdateProjectRequest struct {
	Name            *string `validate:"omitempty,max=128"`
	Transformations struct {
		Add    []CreateTransformationRequest
		Update []UpdateTransformationRequest
		Remove []DeleteTransformationRequest
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
