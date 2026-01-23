package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/apperr"
)

// Image handlers

// CreateUploadURL issues a presigned URL for uploading an image
func (h *handler) CreateUploadURL(ctx echo.Context, projectID ProjectIDPath) error {
	rctx := ctx.Request().Context()

	var req CreateUploadURLRequest
	if err := ctx.Bind(&req); err != nil {
		return apperr.NewError(apperr.CodeBadRequest).
			WithCause(err).
			WithSummary("Failed to parse request body")
	}

	uploadURL, err := h.imageSvc.CreateUploadURL(rctx,
		CreateUploadURLRequestToDomain(projectID, req))
	if err != nil {
		return fmt.Errorf("creating upload url: %w", err)
	}

	return ctx.JSON(http.StatusOK, UploadURLToWeb(uploadURL))
}

// GetImage gets image details
func (h *handler) GetImage(ctx echo.Context, projectID ProjectIDPath, imageID ImageIDPath,
	params GetImageParams,
) error {
	rctx := ctx.Request().Context()

	var (
		image domain.Image
		err   error
	)
	if params.WaitUntilProcessed != nil && *params.WaitUntilProcessed {
		image, err = h.imageSvc.GetWaitUntilProcessed(rctx, imageID)
		if err != nil {
			return fmt.Errorf("getting image with wait until processed: %w", err)
		}
	} else {
		image, err = h.imageSvc.Get(rctx, string(imageID))
		if err != nil {
			return fmt.Errorf("getting image: %w", err)
		}
	}

	return ctx.JSON(http.StatusOK, ImageToWeb(image))
}

// ReprocessImagesAdmin reprocesses multiple images in a project (admin endpoint)
func (h *handler) ReprocessImagesAdmin(ctx echo.Context, projectID ProjectIDPath) error {
	return apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented")
}

// ListImagesAdmin lists all images in a project (admin endpoint)
func (h *handler) ListImagesAdmin(ctx echo.Context, projectID ProjectIDPath,
	params ListImagesAdminParams,
) error {
	rctx := ctx.Request().Context()

	images, err := h.imageSvc.List(rctx, ListImagesAdminParamsToDomain(projectID, params))
	if err != nil {
		return fmt.Errorf("listing images: %w", err)
	}

	return ctx.JSON(http.StatusOK, ImagesToWeb(images))
}

// DeleteImageAdmin deletes an image (admin endpoint)
func (h *handler) DeleteImageAdmin(ctx echo.Context, projectID ProjectIDPath, imageID ImageIDPath,
) error {
	rctx := ctx.Request().Context()

	if err := h.imageSvc.Delete(rctx, imageID); err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}

	return ctx.NoContent(http.StatusOK)
}

// ListImages lists all images in a project
func (h *handler) ListImages(ctx echo.Context, projectID ProjectIDPath, params ListImagesParams,
) error {
	rctx := ctx.Request().Context()

	images, err := h.imageSvc.List(rctx, ListImagesParamsToDomain(projectID, params))
	if err != nil {
		return fmt.Errorf("listing images: %w", err)
	}

	return ctx.JSON(http.StatusOK, ImagesToWeb(images))
}

// DeleteImage deletes an image
func (h *handler) DeleteImage(ctx echo.Context, projectID ProjectIDPath, imageID ImageIDPath,
) error {
	rctx := ctx.Request().Context()

	if err := h.imageSvc.Delete(rctx, imageID); err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}

	return ctx.NoContent(http.StatusOK)
}
