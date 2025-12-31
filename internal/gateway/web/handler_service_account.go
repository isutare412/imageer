package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/pkg/apperr"
)

// Service Account handlers

// GetServiceAccountAdmin gets service account details (admin endpoint)
func (h *handler) GetServiceAccountAdmin(
	ctx echo.Context, serviceAccountID ServiceAccountIDPath,
) error {
	rctx := ctx.Request().Context()

	account, err := h.serviceAccountSvc.GetByID(rctx, serviceAccountID)
	if err != nil {
		return fmt.Errorf("getting service account: %w", err)
	}

	return ctx.JSON(http.StatusOK, ServiceAccountToWeb(account))
}

// ListServiceAccountsAdmin lists service accounts of a project (admin endpoint)
func (h *handler) ListServiceAccountsAdmin(
	ctx echo.Context, params ListServiceAccountsAdminParams,
) error {
	rctx := ctx.Request().Context()

	accounts, err := h.serviceAccountSvc.List(rctx, ListServiceAccountsAdminParamsToDomain(params))
	if err != nil {
		return fmt.Errorf("listing service accounts: %w", err)
	}

	return ctx.JSON(http.StatusOK, ServiceAccountsToWeb(accounts))
}

// CreateServiceAccountAdmin creates a new service account for a project (admin endpoint)
func (h *handler) CreateServiceAccountAdmin(ctx echo.Context) error {
	rctx := ctx.Request().Context()

	var req CreateServiceAccountAdminRequest
	if err := ctx.Bind(&req); err != nil {
		return apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Failed to parse request body").
			WithCause(err)
	}

	account, err := h.serviceAccountSvc.Create(rctx, CreateServiceAccountAdminRequestToDomain(req))
	if err != nil {
		return fmt.Errorf("creating service account: %w", err)
	}

	return ctx.JSON(http.StatusOK, ServiceAccountWithAPIKeyToWeb(account))
}

// UpdateServiceAccountAdmin updates a service account (admin endpoint)
func (h *handler) UpdateServiceAccountAdmin(
	ctx echo.Context, serviceAccountID ServiceAccountIDPath,
) error {
	rctx := ctx.Request().Context()

	var req UpdateServiceAccountAdminRequest
	if err := ctx.Bind(&req); err != nil {
		return apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Failed to parse request body").
			WithCause(err)
	}

	account, err := h.serviceAccountSvc.Update(rctx,
		UpdateServiceAccountAdminRequestToDomain(serviceAccountID, req))
	if err != nil {
		return fmt.Errorf("updating service account: %w", err)
	}

	return ctx.JSON(http.StatusOK, ServiceAccountToWeb(account))
}

// DeleteServiceAccountAdmin deletes a service account (admin endpoint)
func (h *handler) DeleteServiceAccountAdmin(
	ctx echo.Context, serviceAccountID ServiceAccountIDPath,
) error {
	rctx := ctx.Request().Context()

	if err := h.serviceAccountSvc.Delete(rctx, serviceAccountID); err != nil {
		return fmt.Errorf("deleting service account: %w", err)
	}

	return ctx.NoContent(http.StatusOK)
}
