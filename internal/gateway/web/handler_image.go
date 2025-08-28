package web

import (
	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/pkg/apperr"
)

// Image handlers

// CreateUploadURL issues a presigned URL for uploading an image
func (h *handler) CreateUploadURL(ctx echo.Context, projectID ProjectIDPath) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}

// GetImage gets image details
func (h *handler) GetImage(ctx echo.Context, projectID ProjectIDPath, imageID ImageIDPath) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}
