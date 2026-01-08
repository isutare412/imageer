package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/pkg/apperr"
)

// Image handlers

// CreateUploadURL issues a presigned URL for uploading an image
func (h *handler) CreateUploadURL(ctx echo.Context, projectID ProjectIDPath) error {
	rctx := ctx.Request().Context()

	var req CreateUploadURLRequest
	if err := ctx.Bind(&req); err != nil {
		return apperr.NewError(apperr.CodeBadRequest).WithCause(err).WithSummary("Failed to parse request body")
	}

	uploadURL, err := h.imageSvc.CreateUploadURL(rctx,
		CreateUploadURLRequestToDomain(projectID, req))
	if err != nil {
		return fmt.Errorf("creating upload url: %w", err)
	}

	return ctx.JSON(http.StatusOK, UploadURLToWeb(uploadURL))
}

// GetImage gets image details
func (h *handler) GetImage(ctx echo.Context, projectID ProjectIDPath, imageID ImageIDPath) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}
