package domain

import "time"

type Transformation struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Default   bool
	Width     int64
	Height    int64
}

type CreateTransformationRequest struct {
	Name    string `validate:"required,max=64"`
	Default bool
	Width   int64 `validate:"required,min=1"`
	Height  int64 `validate:"required,min=1"`
}

func (r CreateTransformationRequest) ToTransformation() Transformation {
	return Transformation{
		Name:    r.Name,
		Default: r.Default,
		Width:   r.Width,
		Height:  r.Height,
	}
}

type UpsertTransformationRequest struct {
	ID      *string `validate:"omitempty,max=36"`
	Name    *string `validate:"required_without=ID,omitempty,max=64"`
	Default *bool
	Width   *int64 `validate:"required_without=ID,omitempty,min=1"`
	Height  *int64 `validate:"required_without=ID,omitempty,min=1"`
}

func (r UpsertTransformationRequest) IsUpdateRequest() bool {
	return r.ID != nil
}

func (r UpsertTransformationRequest) IsCreateRequest() bool {
	return !r.IsUpdateRequest()
}
