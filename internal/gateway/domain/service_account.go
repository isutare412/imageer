package domain

import (
	"time"

	"github.com/isutare412/imageer/pkg/serviceaccounts"
)

type ServiceAccount struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpireAt  *time.Time
	Name      string
	Authority serviceaccounts.Authority
	APIKey    string
}

type CreateServiceAccountRequest struct {
	Name      string                    `validate:"required,max=128"`
	Authority serviceaccounts.Authority `validate:"required,oneof=FULL_ACCESS PROJECT_ACCESS"`
	ExpireAt  *time.Time                `validate:"omitempty,gt"`
}

type UpdateServiceAccountRequest struct {
	Name      *string                    `validate:"omitempty,max=128"`
	Authority *serviceaccounts.Authority `validate:"omitempty,oneof=FULL_ACCESS PROJECT_ACCESS"`
	ExpireAt  *time.Time                 `validate:"omitempty,gt"`
}

type ServiceAccounts struct {
	Items []ServiceAccount
	Total int64
}
