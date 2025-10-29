package web

import (
	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/pkg/apperr"
)

// Service Account handlers

// ListServiceAccountsAdmin lists service accounts of a project (admin endpoint)
func (h *handler) ListServiceAccountsAdmin(ctx echo.Context, params ListServiceAccountsAdminParams) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}

// CreateServiceAccountAdmin creates a new service account for a project (admin endpoint)
func (h *handler) CreateServiceAccountAdmin(ctx echo.Context) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}

// DeleteServiceAccountAdmin deletes a service account (admin endpoint)
func (h *handler) DeleteServiceAccountAdmin(ctx echo.Context, serviceAccountID ServiceAccountIDPath) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}

// GetServiceAccountAdmin gets service account details (admin endpoint)
func (h *handler) GetServiceAccountAdmin(ctx echo.Context, serviceAccountID ServiceAccountIDPath) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}

// UpdateServiceAccountAdmin updates a service account (admin endpoint)
func (h *handler) UpdateServiceAccountAdmin(ctx echo.Context, serviceAccountID ServiceAccountIDPath) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}
