package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/pkg/images"
)

type Image struct {
	ID        string `gorm:"size:36"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FileName  string        `gorm:"size:1024"`
	Format    images.Format `gorm:"size:32"`
	State     images.State  `gorm:"size:32"`

	ProjectID string   `gorm:"size:36"`
	Project   *Project `gorm:"constraint:OnDelete:CASCADE"`
}

func (i *Image) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.NewString()
	}
	return nil
}
