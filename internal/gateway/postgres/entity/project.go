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

	Transformations []Transformation `gorm:"constraint:OnDelete:CASCADE"`
}

func NewProject(req domain.Project) Project {
	return Project{
		Name: req.Name,
		Transformations: lo.Map(req.Transformations, func(t domain.Transformation, _ int) Transformation {
			return NewTransformation(t)
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
		Transformations: lo.Map(p.Transformations, func(t Transformation, _ int) domain.Transformation {
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
