package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/isutare412/imageer/pkg/serviceaccounts"
)

type ServiceAccount struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpireAt  *time.Time
	Name      string
	Authority serviceaccounts.Authority
	Token     string
}

type CreateServiceAccountRequest struct {
	Name      string                    `validate:"required,max=128"`
	Authority serviceaccounts.Authority `validate:"required,max=24"`
	ExpireAt  *time.Time                `validate:"omitempty,gt"`
}

type UpdateServiceAccountRequest struct {
	Name      *string                    `validate:"omitempty,max=128"`
	Authority *serviceaccounts.Authority `validate:"omitempty,required,max=24"`
	ExpireAt  *time.Time                 `validate:"omitempty,gt"`
}

type ServiceAccounts struct {
	Items []ServiceAccount
	Total int64
}
