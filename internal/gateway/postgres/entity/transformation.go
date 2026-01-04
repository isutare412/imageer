package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/images"
)

type Transformation struct {
	ID        string `gorm:"size:36"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:64; uniqueIndex:idx_project_id_name,priority:2"`
	Default   bool

	Format  images.Format `gorm:"size:32"`
	Quality images.Quality
	Fit     *images.Fit `gorm:"size:32"`
	Width   *int64
	Height  *int64

	Crop   bool
	Anchor *images.Anchor `gorm:"size:32"`

	ProjectID string `gorm:"size:36; uniqueIndex:idx_project_id_name,priority:1"`
}

func NewTransformation(t domain.Transformation) Transformation {
	return Transformation{
		Name:    t.Name,
		Default: t.Default,
		Format:  t.Format,
		Quality: t.Quality,
		Fit:     t.Fit,
		Width:   t.Width,
		Height:  t.Height,
		Crop:    t.Crop,
		Anchor:  t.Anchor,
	}
}

func NewTransformationFromUpsert(
	projID string, req domain.UpsertTransformationRequest,
) Transformation {
	return Transformation{
		Name:      lo.FromPtr(req.Name),
		Default:   lo.FromPtr(req.Default),
		Format:    req.Format.GetOrDefault(),
		Quality:   req.Quality.GetOrDefault(),
		Fit:       req.Fit,
		Width:     req.Width,
		Height:    req.Height,
		Crop:      lo.FromPtr(req.Crop),
		Anchor:    req.Anchor,
		ProjectID: projID,
	}
}

func (t *Transformation) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return nil
}

func (t Transformation) ToDomain() domain.Transformation {
	return domain.Transformation{
		ID:        t.ID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		Name:      t.Name,
		Default:   t.Default,
		Format:    t.Format,
		Quality:   t.Quality,
		Fit:       t.Fit,
		Width:     t.Width,
		Height:    t.Height,
		Crop:      t.Crop,
		Anchor:    t.Anchor,
	}
}
