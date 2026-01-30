package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
	"github.com/isutare412/imageer/pkg/apperr"
)

// Service Account handlers

// GetServiceAccountAdmin gets service account details (admin endpoint)
func (h *Handler) GetServiceAccountAdmin(
	w http.ResponseWriter, r *http.Request, serviceAccountID gen.ServiceAccountIDPath,
) {
	ctx := r.Context()

	account, err := h.serviceAccountSvc.GetByID(ctx, serviceAccountID)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("getting service account: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ServiceAccountToWeb(account))
}

// ListServiceAccountsAdmin lists service accounts of a project (admin endpoint)
func (h *Handler) ListServiceAccountsAdmin(
	w http.ResponseWriter, r *http.Request, params gen.ListServiceAccountsAdminParams,
) {
	ctx := r.Context()

	accounts, err := h.serviceAccountSvc.List(ctx, ListServiceAccountsAdminParamsToDomain(params))
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("listing service accounts: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ServiceAccountsToWeb(accounts))
}

// CreateServiceAccountAdmin creates a new service account for a project (admin endpoint)
func (h *Handler) CreateServiceAccountAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req gen.CreateServiceAccountAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		gen.RespondError(w, r, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Failed to parse request body").
			WithCause(err))
		return
	}

	account, err := h.serviceAccountSvc.Create(ctx, CreateServiceAccountAdminRequestToDomain(req))
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("creating service account: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ServiceAccountWithAPIKeyToWeb(account))
}

// UpdateServiceAccountAdmin updates a service account (admin endpoint)
func (h *Handler) UpdateServiceAccountAdmin(
	w http.ResponseWriter, r *http.Request, serviceAccountID gen.ServiceAccountIDPath,
) {
	ctx := r.Context()

	var req gen.UpdateServiceAccountAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		gen.RespondError(w, r, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Failed to parse request body").
			WithCause(err))
		return
	}

	account, err := h.serviceAccountSvc.Update(ctx,
		UpdateServiceAccountAdminRequestToDomain(serviceAccountID, req))
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("updating service account: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ServiceAccountToWeb(account))
}

// DeleteServiceAccountAdmin deletes a service account (admin endpoint)
func (h *Handler) DeleteServiceAccountAdmin(
	w http.ResponseWriter, r *http.Request, serviceAccountID gen.ServiceAccountIDPath,
) {
	ctx := r.Context()

	if err := h.serviceAccountSvc.Delete(ctx, serviceAccountID); err != nil {
		gen.RespondError(w, r, fmt.Errorf("deleting service account: %w", err))
		return
	}

	gen.RespondNoContent(w, http.StatusOK)
}
