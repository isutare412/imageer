package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/images"
)

type Preset struct {
	ID        string `gorm:"size:36"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:64; uniqueIndex:idx_project_id_name,priority:2"`
	Default   bool

	Format  images.Format  `gorm:"size:32"`
	Quality images.Quality `gorm:"type:smallint"`
	Fit     *images.Fit    `gorm:"size:32"`
	Anchor  *images.Anchor `gorm:"size:32"`
	Width   *int64         `gorm:"type:integer"`
	Height  *int64         `gorm:"type:integer"`

	ProjectID string `gorm:"size:36; uniqueIndex:idx_project_id_name,priority:1"`
}

func NewPreset(t domain.Preset) Preset {
	return Preset{
		Name:    t.Name,
		Default: t.Default,
		Format:  t.Format,
		Quality: t.Quality,
		Fit:     t.Fit,
		Anchor:  t.Anchor,
		Width:   t.Width,
		Height:  t.Height,
	}
}

func NewPresetFromUpsert(
	projID string, req domain.UpsertPresetRequest,
) Preset {
	return Preset{
		Name:      lo.FromPtr(req.Name),
		Default:   lo.FromPtr(req.Default),
		Format:    req.Format.GetOrDefault(),
		Quality:   req.Quality.GetOrDefault(),
		Fit:       req.Fit,
		Anchor:    req.Anchor,
		Width:     req.Width,
		Height:    req.Height,
		ProjectID: projID,
	}
}

func (t *Preset) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return nil
}

func (t Preset) ToDomain() domain.Preset {
	return domain.Preset{
		ID:        t.ID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		Name:      t.Name,
		Default:   t.Default,
		Format:    t.Format,
		Quality:   t.Quality,
		Fit:       t.Fit,
		Anchor:    t.Anchor,
		Width:     t.Width,
		Height:    t.Height,
	}
}
