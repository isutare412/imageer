package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/images"
)

type ImageVariant struct {
	ID        string `gorm:"size:36"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Format    images.Format       `gorm:"size:32"`
	State     images.VariantState `gorm:"size:32"`
	S3Key     string              `gorm:"size:1024"`
	URL       string              `gorm:"size:1024"`

	ImageID  string `gorm:"size:36; uniqueIndex:idx_image_id_preset_id,priority:1"`
	PresetID string `gorm:"size:36; uniqueIndex:idx_image_id_preset_id,priority:2; index"`
	Preset   Preset `gorm:"constraint:OnDelete:SET NULL"`
}

func NewImageVariant(iv domain.ImageVariant) ImageVariant {
	return ImageVariant{
		Format:   iv.Format,
		State:    iv.State,
		S3Key:    iv.S3Key,
		URL:      iv.URL,
		ImageID:  iv.ImageID,
		PresetID: iv.Preset.ID,
	}
}

func (i *ImageVariant) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.NewString()
	}
	return nil
}

func (i ImageVariant) ToDomain() domain.ImageVariant {
	return domain.ImageVariant{
		ID:        i.ID,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
		Format:    i.Format,
		State:     i.State,
		S3Key:     i.S3Key,
		URL:       i.URL,
		ImageID:   i.ImageID,
		Preset:    i.Preset.ToReference(),
	}
}
