package domain

import (
	"time"

	"github.com/isutare412/imageer/pkg/serviceaccounts"
)

type ServiceAccount struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpireAt    *time.Time
	Name        string
	AccessScope serviceaccounts.AccessScope
	Projects    []ProjectReference
	APIKey      string
}

type CreateServiceAccountRequest struct {
	Name        string                      `validate:"required,max=128"`
	AccessScope serviceaccounts.AccessScope `validate:"required,oneof=FULL PROJECT"`
	ProjectIDs  []string                    `validate:"dive,required"`
	ExpireAt    *time.Time                  `validate:"omitempty,gt"`
}

type UpdateServiceAccountRequest struct {
	ID          string                       `validate:"max=36"`
	Name        *string                      `validate:"omitempty,max=128"`
	AccessScope *serviceaccounts.AccessScope `validate:"omitempty,oneof=FULL PROJECT"`
	ProjectIDs  []string                     `validate:"dive,required"`
	ExpireAt    *time.Time
}

type ServiceAccounts struct {
	Items []ServiceAccount
	Total int64
}
