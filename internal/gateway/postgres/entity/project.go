package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/pkg/apperr"
)

type Project struct {
	ID        string `gorm:"size:36"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:128"`

	Transformations []*Transformation `gorm:"constraint:OnDelete:CASCADE"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return nil
}
