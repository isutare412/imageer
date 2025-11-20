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
	Width   int64 `validate:"min=1"`
	Height  int64 `validate:"min=1"`
}

type UpdateTransformationRequest struct {
	ID      string  `validate:"required"`
	Name    *string `validate:"omitempty,max=64"`
	Default *bool
	Width   *int64 `validate:"omitempty,min=1"`
	Height  *int64 `validate:"omitempty,min=1"`
}

type DeleteTransformationRequest struct {
	ID string `validate:"required"`
}
