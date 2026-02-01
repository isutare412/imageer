package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/trace"
)

// Project handlers

// GetProject gets project details
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request, projectID gen.ProjectIDPath) {
	ctx, span := trace.StartSpan(r.Context(), "web.handlers.GetProject")
	defer span.End()

	project, err := h.projectSvc.GetByID(ctx, projectID)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("getting project by id: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ProjectToWeb(project))
}

// GetProjectAdmin gets project details (admin endpoint)
func (h *Handler) GetProjectAdmin(w http.ResponseWriter, r *http.Request, projectID gen.ProjectIDPath) {
	ctx, span := trace.StartSpan(r.Context(), "web.handlers.GetProjectAdmin")
	defer span.End()

	project, err := h.projectSvc.GetByID(ctx, projectID)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("getting project by id: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ProjectToWeb(project))
}

// Admin Project handlers

// ListProjectsAdmin lists all projects (admin endpoint)
func (h *Handler) ListProjectsAdmin(w http.ResponseWriter, r *http.Request, params gen.ListProjectsAdminParams) {
	ctx, span := trace.StartSpan(r.Context(), "web.handlers.ListProjectsAdmin")
	defer span.End()

	projects, err := h.projectSvc.List(ctx, ListProjectsAdminParamsToDomain(params))
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("listing projects: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ProjectsToWeb(projects))
}

// CreateProjectAdmin creates a new project (admin endpoint)
func (h *Handler) CreateProjectAdmin(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "web.handlers.CreateProjectAdmin")
	defer span.End()

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
func (h *Handler) UpdateProjectAdmin(w http.ResponseWriter, r *http.Request, projectID gen.ProjectIDPath) {
	ctx, span := trace.StartSpan(r.Context(), "web.handlers.UpdateProjectAdmin")
	defer span.End()

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
func (h *Handler) DeleteProjectAdmin(w http.ResponseWriter, r *http.Request, projectID gen.ProjectIDPath) {
	ctx, span := trace.StartSpan(r.Context(), "web.handlers.DeleteProjectAdmin")
	defer span.End()

	if err := h.projectSvc.Delete(ctx, projectID); err != nil {
		gen.RespondError(w, r, fmt.Errorf("deleting project: %w", err))
		return
	}

	gen.RespondNoContent(w, http.StatusOK)
}
