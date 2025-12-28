package domain

import (
	"time"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/pkg/dbhelpers"
	"github.com/isutare412/imageer/pkg/serviceaccounts"
)

type ServiceAccount struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpireAt    *time.Time
	Name        string
	AccessScope serviceaccounts.AccessScope
	APIKey      string
	Projects    []ProjectReference
}

type CreateServiceAccountRequest struct {
	Name        string                      `validate:"required,max=128"`
	AccessScope serviceaccounts.AccessScope `validate:"required,validateFn=IsAAccessScope"`
	ProjectIDs  []string                    `validate:"dive,required"`
	ExpireAt    *time.Time                  `validate:"omitempty,gt"`
}

type UpdateServiceAccountRequest struct {
	ID          string                       `validate:"max=36"`
	Name        *string                      `validate:"omitempty,max=128"`
	AccessScope *serviceaccounts.AccessScope `validate:"omitempty,validateFn=IsAAccessScope"`
	ProjectIDs  []string                     `validate:"dive,required"`
	ExpireAt    *time.Time
}

type ServiceAccounts struct {
	Items []ServiceAccount
	Total int64
}

type ListServiceAccountsParams struct {
	Offset *int `validate:"omitempty,min=0"`
	Limit  *int `validate:"omitempty,min=1,max=100"`

	SearchFilter ServiceAccountSearchFilter
	SortFilter   ServiceAccountSortFilter
}

func (p ListServiceAccountsParams) OffsetOrDefault() int {
	return lo.FromPtrOr(p.Offset, 0)
}

func (p ListServiceAccountsParams) LimitOrDefault() int {
	return lo.FromPtrOr(p.Limit, 20)
}

type ServiceAccountSearchFilter struct {
	Name *string
}

type ServiceAccountSortFilter struct {
	CreatedAt bool
	UpdatedAt bool
	Direction dbhelpers.SortDirection
}
