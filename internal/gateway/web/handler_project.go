package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/pkg/apperr"
)

// Project handlers

// GetProject gets project details
func (h *handler) GetProject(ctx echo.Context, projectID ProjectIDPath) error {
	rctx := ctx.Request().Context()

	project, err := h.projectSvc.GetByID(rctx, projectID)
	if err != nil {
		return fmt.Errorf("getting project by id: %w", err)
	}

	return ctx.JSON(http.StatusOK, ProjectToWeb(project))
}

// GetProjectAdmin gets project details (admin endpoint)
func (h *handler) GetProjectAdmin(ctx echo.Context, projectID ProjectIDPath) error {
	rctx := ctx.Request().Context()

	project, err := h.projectSvc.GetByID(rctx, projectID)
	if err != nil {
		return fmt.Errorf("getting project by id: %w", err)
	}

	return ctx.JSON(http.StatusOK, ProjectToWeb(project))
}

// Admin Project handlers

// ListProjectsAdmin lists all projects (admin endpoint)
func (h *handler) ListProjectsAdmin(ctx echo.Context, params ListProjectsAdminParams) error {
	rctx := ctx.Request().Context()

	projects, err := h.projectSvc.List(rctx, ListProjectsAdminParamsToDomain(params))
	if err != nil {
		return fmt.Errorf("listing projects: %w", err)
	}

	return ctx.JSON(http.StatusOK, ProjectsToWeb(projects))
}

// CreateProjectAdmin creates a new project (admin endpoint)
func (h *handler) CreateProjectAdmin(ctx echo.Context) error {
	rctx := ctx.Request().Context()

	var req CreateProjectAdminRequest
	if err := ctx.Bind(&req); err != nil {
		return apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Failed to parse request body").
			WithCause(err)
	}

	project, err := h.projectSvc.Create(rctx, CreateProjectAdminRequestToDomain(req))
	if err != nil {
		return fmt.Errorf("creating project: %w", err)
	}

	return ctx.JSON(http.StatusOK, ProjectToWeb(project))
}

// UpdateProjectAdmin updates project details (admin endpoint)
func (h *handler) UpdateProjectAdmin(ctx echo.Context, projectID ProjectIDPath) error {
	rctx := ctx.Request().Context()

	var req UpdateProjectAdminRequest
	if err := ctx.Bind(&req); err != nil {
		return apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Failed to parse request body").
			WithCause(err)
	}

	project, err := h.projectSvc.Update(rctx, UpdateProjectAdminRequestToDomain(projectID, req))
	if err != nil {
		return fmt.Errorf("updating project: %w", err)
	}

	return ctx.JSON(http.StatusOK, ProjectToWeb(project))
}

func (h *handler) DeleteProjectAdmin(ctx echo.Context, projectID ProjectIDPath) error {
	rctx := ctx.Request().Context()

	if err := h.projectSvc.Delete(rctx, projectID); err != nil {
		return fmt.Errorf("deleting project: %w", err)
	}
	return ctx.NoContent(http.StatusOK)
}

// ReprocessImagesAdmin reprocesses multiple images in a project (admin endpoint)
func (h *handler) ReprocessImagesAdmin(ctx echo.Context, projectID ProjectIDPath) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}
