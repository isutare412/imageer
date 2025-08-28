package web

import (
	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/pkg/apperr"
)

// Project handlers

// GetProject gets project details
func (h *handler) GetProject(ctx echo.Context, projectID ProjectIDPath) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}

// Admin Project handlers

// ListProjectsAdmin lists all projects (admin endpoint)
func (h *handler) ListProjectsAdmin(ctx echo.Context, params ListProjectsAdminParams) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}

// CreateProjectAdmin creates a new project (admin endpoint)
func (h *handler) CreateProjectAdmin(ctx echo.Context) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}

// GetProjectAdmin gets project details (admin endpoint)
func (h *handler) GetProjectAdmin(ctx echo.Context, projectID ProjectIDPath) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}

// UpdateProjectAdmin updates project details (admin endpoint)
func (h *handler) UpdateProjectAdmin(ctx echo.Context, projectID ProjectIDPath) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}

// ReprocessImagesAdmin reprocesses multiple images in a project (admin endpoint)
func (h *handler) ReprocessImagesAdmin(ctx echo.Context, projectID ProjectIDPath) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}
