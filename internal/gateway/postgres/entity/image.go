package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/images"
)

type Image struct {
	ID               string `gorm:"size:36"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	OriginalFileName string             `gorm:"size:256"`
	ContentType      images.ContentType `gorm:"size:64"`
	State            images.State       `gorm:"size:32"`

	ProjectID string   `gorm:"size:36"`
	Project   *Project `gorm:"constraint:OnDelete:CASCADE"`
}

func (i *Image) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return apperr.NewError(apperr.CodeInternalServerError).
				WithSummary("failed to generate UUIDv7 for image ID").
				WithCause(err)
		}

		i.ID = id.String()
	}
	return nil
}
