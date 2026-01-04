package domain

import (
	"time"

	"github.com/isutare412/imageer/pkg/images"
)

type Transformation struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Default   bool

	Format  images.Format
	Quality images.Quality
	Fit     *images.Fit
	Width   *int64
	Height  *int64

	Crop   bool
	Anchor *images.Anchor
}

type CreateTransformationRequest struct {
	Name    string `validate:"required,max=64"`
	Default bool

	Format  *images.Format  `validate:"omitempty,validateFn=ValidateForTransformation"`
	Quality *images.Quality `validate:"omitempty,validateFn=Validate"`
	Fit     *images.Fit     `validate:"omitempty,validateFn=Validate"`
	Width   *int64          `validate:"omitempty,min=1,max=4000"`
	Height  *int64          `validate:"omitempty,min=1,max=4000"`

	Crop   bool
	Anchor *images.Anchor `validate:"omitempty,validateFn=Validate"`
}

func (r CreateTransformationRequest) ToTransformation() Transformation {
	return Transformation{
		Name:    r.Name,
		Default: r.Default,
		Format:  r.Format.GetOrDefault(),
		Quality: r.Quality.GetOrDefault(),
		Fit:     r.Fit,
		Width:   r.Width,
		Height:  r.Height,
		Anchor:  r.Anchor,
	}
}

type UpsertTransformationRequest struct {
	ID      *string `validate:"omitempty,max=36"`
	Name    *string `validate:"omitempty,max=64"`
	Default *bool

	Format  *images.Format  `validate:"omitempty,validateFn=ValidateForTransformation"`
	Quality *images.Quality `validate:"omitempty,validateFn=Validate"`
	Fit     *images.Fit     `validate:"omitempty,validateFn=Validate"`
	Width   *int64          `validate:"omitempty,min=1,max=4000"`
	Height  *int64          `validate:"omitempty,min=1,max=4000"`

	Crop   *bool
	Anchor *images.Anchor `validate:"omitempty,validateFn=Validate"`
}

func (r UpsertTransformationRequest) IsUpdateRequest() bool {
	return r.ID != nil
}

func (r UpsertTransformationRequest) IsCreateRequest() bool {
	return !r.IsUpdateRequest()
}
