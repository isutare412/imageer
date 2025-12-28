package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

type Transformation struct {
	ID        string `gorm:"size:36"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:64; uniqueIndex:idx_project_id_name,priority:2"`
	Default   bool
	Width     int64
	Height    int64

	ProjectID string `gorm:"size:36; uniqueIndex:idx_project_id_name,priority:1"`
}

func NewTransformation(projID string, req domain.UpsertTransformationRequest) Transformation {
	return Transformation{
		Name:      lo.FromPtr(req.Name),
		Default:   lo.FromPtr(req.Default),
		Width:     lo.FromPtr(req.Width),
		Height:    lo.FromPtr(req.Height),
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
		Width:     t.Width,
		Height:    t.Height,
	}
}
