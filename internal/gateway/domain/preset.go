package domain

import (
	"time"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/pkg/dbhelpers"
	"github.com/isutare412/imageer/pkg/images"
)

type Preset struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Default   bool

	Format  images.Format
	Quality images.Quality
	Fit     *images.Fit
	Anchor  *images.Anchor
	Width   *int64
	Height  *int64
}

type PresetReference struct {
	ID   string
	Name string
}

type CreatePresetRequest struct {
	Name    string `validate:"required,max=64,kebabcase"`
	Default bool

	Format  *images.Format  `validate:"omitempty,validateFn=ValidateForPreset"`
	Quality *images.Quality `validate:"omitempty,validateFn=Validate"`
	Fit     *images.Fit     `validate:"omitempty,validateFn=Validate"`
	Anchor  *images.Anchor  `validate:"omitempty,validateFn=Validate"`
	Width   *int64          `validate:"omitempty,min=1,max=4000"`
	Height  *int64          `validate:"omitempty,min=1,max=4000"`
}

func (r CreatePresetRequest) ToPreset() Preset {
	return Preset{
		Name:    r.Name,
		Default: r.Default,
		Format:  r.Format.GetOrDefault(),
		Quality: r.Quality.GetOrDefault(),
		Fit:     r.Fit,
		Anchor:  r.Anchor,
		Width:   r.Width,
		Height:  r.Height,
	}
}

type UpsertPresetRequest struct {
	ID      *string `validate:"omitempty,max=36"`
	Name    *string `validate:"omitempty,max=64,kebabcase"`
	Default *bool

	Format  *images.Format  `validate:"omitempty,validateFn=ValidateForPreset"`
	Quality *images.Quality `validate:"omitempty,validateFn=Validate"`
	Fit     *images.Fit     `validate:"omitempty,validateFn=Validate"`
	Anchor  *images.Anchor  `validate:"omitempty,validateFn=Validate"`
	Width   *int64          `validate:"omitempty,min=1,max=4000"`
	Height  *int64          `validate:"omitempty,min=1,max=4000"`
}

func (r UpsertPresetRequest) IsUpdateRequest() bool {
	return r.ID != nil
}

func (r UpsertPresetRequest) IsCreateRequest() bool {
	return !r.IsUpdateRequest()
}

type ListPresetsParams struct {
	Offset *int `validate:"omitempty,min=0"`
	Limit  *int `validate:"omitempty,min=1,max=100"`

	SearchFilter PresetSearchFilter
	SortFilter   PresetSortFilter
}

func (p ListPresetsParams) OffsetOrDefault() int {
	return lo.FromPtrOr(p.Offset, -1)
}

func (p ListPresetsParams) LimitOrDefault() int {
	return lo.FromPtrOr(p.Limit, -1)
}

type PresetSearchFilter struct {
	ProjectID *string
	Names     []string
}

type PresetSortFilter struct {
	CreatedAt bool
	UpdatedAt bool
	Direction dbhelpers.SortDirection
}
