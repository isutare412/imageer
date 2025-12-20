package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/serviceaccounts"
)

type ServiceAccount struct {
	ID          string `gorm:"size:36"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string                      `gorm:"size:128"`
	AccessScope serviceaccounts.AccessScope `gorm:"size:32"`
	ExpireAt    time.Time
	APIKey      string `gorm:"size:64; uniqueIndex"`

	Projects []*Project `gorm:"many2many:service_account_projects"`
}

func (sa *ServiceAccount) BeforeCreate(tx *gorm.DB) error {
	if sa.ID == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return apperr.NewError(apperr.CodeInternalServerError).
				WithSummary("failed to generate UUIDv7 for service account ID").
				WithCause(err)
		}
		sa.ID = id.String()
	}
	return nil
}
