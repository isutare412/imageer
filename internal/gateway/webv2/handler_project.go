package webv2

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
	"github.com/isutare412/imageer/pkg/apperr"
)

// Project handlers

// GetProject gets project details
func (h *handler) GetProject(w http.ResponseWriter, r *http.Request, projectID gen.ProjectIDPath) {
	ctx := r.Context()

	project, err := h.projectSvc.GetByID(ctx, projectID)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("getting project by id: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ProjectToWeb(project))
}

// GetProjectAdmin gets project details (admin endpoint)
func (h *handler) GetProjectAdmin(w http.ResponseWriter, r *http.Request, projectID gen.ProjectIDPath) {
	ctx := r.Context()

	project, err := h.projectSvc.GetByID(ctx, projectID)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("getting project by id: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ProjectToWeb(project))
}

// Admin Project handlers

// ListProjectsAdmin lists all projects (admin endpoint)
func (h *handler) ListProjectsAdmin(w http.ResponseWriter, r *http.Request, params gen.ListProjectsAdminParams) {
	ctx := r.Context()

	projects, err := h.projectSvc.List(ctx, ListProjectsAdminParamsToDomain(params))
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("listing projects: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ProjectsToWeb(projects))
}

// CreateProjectAdmin creates a new project (admin endpoint)
func (h *handler) CreateProjectAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req gen.CreateProjectAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		gen.RespondError(w, r, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Failed to parse request body").
			WithCause(err))
		return
	}

	project, err := h.projectSvc.Create(ctx, CreateProjectAdminRequestToDomain(req))
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("creating project: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ProjectToWeb(project))
}

// UpdateProjectAdmin updates project details (admin endpoint)
func (h *handler) UpdateProjectAdmin(w http.ResponseWriter, r *http.Request, projectID gen.ProjectIDPath) {
	ctx := r.Context()

	var req gen.UpdateProjectAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		gen.RespondError(w, r, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Failed to parse request body").
			WithCause(err))
		return
	}

	project, err := h.projectSvc.Update(ctx, UpdateProjectAdminRequestToDomain(projectID, req))
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("updating project: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ProjectToWeb(project))
}

// DeleteProjectAdmin deletes a project (admin endpoint)
func (h *handler) DeleteProjectAdmin(w http.ResponseWriter, r *http.Request, projectID gen.ProjectIDPath) {
	ctx := r.Context()

	if err := h.projectSvc.Delete(ctx, projectID); err != nil {
		gen.RespondError(w, r, fmt.Errorf("deleting project: %w", err))
		return
	}

	gen.RespondNoContent(w, http.StatusOK)
}
