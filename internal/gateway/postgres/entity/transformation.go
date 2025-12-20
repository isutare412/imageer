package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/pkg/apperr"
)

type Transformation struct {
	ID        string `gorm:"size:36"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:128; uniqueIndex:idx_project_id_name,priority:2"`
	Default   bool
	Width     int64
	Height    int64

	ProjectID string `gorm:"size:36; uniqueIndex:idx_project_id_name,priority:1"`
}

func (t *Transformation) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return apperr.NewError(apperr.CodeInternalServerError).
				WithSummary("failed to generate UUIDv7 for transformation ID").
				WithCause(err)
		}

		t.ID = id.String()
	}
	return nil
}
