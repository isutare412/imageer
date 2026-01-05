package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

type Project struct {
	ID        string `gorm:"size:36"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:128"`

	Presets []Preset `gorm:"constraint:OnDelete:CASCADE"`
}

func NewProject(req domain.Project) Project {
	return Project{
		Name: req.Name,
		Presets: lo.Map(req.Presets, func(t domain.Preset, _ int) Preset {
			return NewPreset(t)
		}),
	}
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return nil
}

func (p Project) ToDomain() domain.Project {
	return domain.Project{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Name:      p.Name,
		Presets: lo.Map(p.Presets, func(t Preset, _ int) domain.Preset {
			return t.ToDomain()
		}),
	}
}

func (p Project) ToReference() domain.ProjectReference {
	return domain.ProjectReference{
		ID:   p.ID,
		Name: p.Name,
	}
}
