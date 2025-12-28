package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/serviceaccounts"
)

type ServiceAccount struct {
	ID          string `gorm:"size:36"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string                      `gorm:"size:128"`
	AccessScope serviceaccounts.AccessScope `gorm:"size:32"`
	ExpireAt    *time.Time
	APIKey      string `gorm:"size:64; uniqueIndex"`

	Projects []*Project `gorm:"many2many:service_account_projects"`
}

func NewServiceAccount(acc domain.ServiceAccount) ServiceAccount {
	return ServiceAccount{
		Name:        acc.Name,
		AccessScope: acc.AccessScope,
		ExpireAt:    acc.ExpireAt,
		APIKey:      acc.APIKey,
	}
}

func (sa *ServiceAccount) BeforeCreate(tx *gorm.DB) error {
	if sa.ID == "" {
		sa.ID = uuid.NewString()
	}
	return nil
}

func (sa ServiceAccount) ToDomain() domain.ServiceAccount {
	return domain.ServiceAccount{
		ID:          sa.ID,
		CreatedAt:   sa.CreatedAt,
		UpdatedAt:   sa.UpdatedAt,
		Name:        sa.Name,
		AccessScope: sa.AccessScope,
		ExpireAt:    sa.ExpireAt,
		APIKey:      sa.APIKey,
		Projects: lo.Map(sa.Projects, func(p *Project, _ int) domain.ProjectReference {
			return p.ToReference()
		}),
	}
}

type ServiceAccountProject struct {
	ServiceAccountID string `gorm:"size:36; primaryKey"`
	ProjectID        string `gorm:"size:36; primaryKey"`
}
