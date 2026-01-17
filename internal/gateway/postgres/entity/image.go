package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/images"
)

type Image struct {
	ID        string `gorm:"size:36"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FileName  string        `gorm:"size:512"`
	Format    images.Format `gorm:"size:32"`
	State     images.State  `gorm:"size:32"`
	S3Key     string        `gorm:"size:1024"`
	URL       string        `gorm:"size:1024"`

	ProjectID string  `gorm:"size:36; index"`
	Project   Project `gorm:"constraint:OnDelete:SET NULL"`

	Variants []ImageVariant `gorm:"constraint:OnDelete:SET NULL"`
}

func NewImage(img domain.Image) Image {
	return Image{
		ID:        img.ID,
		FileName:  img.FileName,
		Format:    img.Format,
		State:     img.State,
		S3Key:     img.S3Key,
		URL:       img.URL,
		ProjectID: img.Project.ID,
	}
}

func (i *Image) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.NewString()
	}
	return nil
}

func (i Image) ToDomain() domain.Image {
	return domain.Image{
		ID:        i.ID,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
		FileName:  i.FileName,
		Format:    i.Format,
		State:     i.State,
		S3Key:     i.S3Key,
		URL:       i.URL,
		Project:   i.Project.ToReference(),
		Variants: lo.Map(i.Variants, func(v ImageVariant, _ int) domain.ImageVariant {
			return v.ToDomain()
		}),
	}
}
