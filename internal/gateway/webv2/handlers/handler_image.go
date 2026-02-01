package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/tracing"
)

// Image handlers

// CreateUploadURL issues a presigned URL for uploading an image
func (h *Handler) CreateUploadURL(w http.ResponseWriter, r *http.Request,
	projectID gen.ProjectIDPath,
) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.CreateUploadURL")
	defer span.End()

	var req gen.CreateUploadURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		gen.RespondError(w, r, apperr.NewError(apperr.CodeBadRequest).
			WithCause(err).
			WithSummary("Failed to parse request body"))
		return
	}

	uploadURL, err := h.imageSvc.CreateUploadURL(ctx,
		CreateUploadURLRequestToDomain(projectID, req))
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("creating upload url: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, UploadURLToWeb(uploadURL))
}

// GetImage gets image details
func (h *Handler) GetImage(
	w http.ResponseWriter, r *http.Request, projectID gen.ProjectIDPath, imageID gen.ImageIDPath,
	params gen.GetImageParams,
) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.GetImage")
	defer span.End()

	var (
		image domain.Image
		err   error
	)
	if params.WaitUntilProcessed != nil && *params.WaitUntilProcessed {
		image, err = h.imageSvc.GetWaitUntilProcessed(ctx, imageID)
		if err != nil {
			gen.RespondError(w, r, fmt.Errorf("getting image with wait until processed: %w", err))
			return
		}
	} else {
		image, err = h.imageSvc.Get(ctx, string(imageID))
		if err != nil {
			gen.RespondError(w, r, fmt.Errorf("getting image: %w", err))
			return
		}
	}

	gen.RespondJSON(w, http.StatusOK, ImageToWeb(image))
}

// ReprocessImagesAdmin reprocesses multiple images in a project (admin endpoint)
func (h *Handler) ReprocessImagesAdmin(w http.ResponseWriter, r *http.Request,
	projectID gen.ProjectIDPath,
) {
	gen.RespondError(w, r, apperr.NewError(apperr.CodeNotImplemented).
		WithSummary("Method not implemented"))
}

// ListImagesAdmin lists all images in a project (admin endpoint)
func (h *Handler) ListImagesAdmin(
	w http.ResponseWriter, r *http.Request, projectID gen.ProjectIDPath,
	params gen.ListImagesAdminParams,
) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.ListImagesAdmin")
	defer span.End()

	images, err := h.imageSvc.List(ctx, ListImagesAdminParamsToDomain(projectID, params))
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("listing images: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ImagesToWeb(images))
}

// DeleteImageAdmin deletes an image (admin endpoint)
func (h *Handler) DeleteImageAdmin(
	w http.ResponseWriter, r *http.Request, projectID gen.ProjectIDPath, imageID gen.ImageIDPath,
) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.DeleteImageAdmin")
	defer span.End()

	if err := h.imageSvc.Delete(ctx, imageID); err != nil {
		gen.RespondError(w, r, fmt.Errorf("deleting image: %w", err))
		return
	}

	gen.RespondNoContent(w, http.StatusOK)
}

// ListImages lists all images in a project
func (h *Handler) ListImages(
	w http.ResponseWriter, r *http.Request, projectID gen.ProjectIDPath,
	params gen.ListImagesParams,
) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.ListImages")
	defer span.End()

	images, err := h.imageSvc.List(ctx, ListImagesParamsToDomain(projectID, params))
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("listing images: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ImagesToWeb(images))
}

// DeleteImage deletes an image
func (h *Handler) DeleteImage(
	w http.ResponseWriter, r *http.Request, projectID gen.ProjectIDPath, imageID gen.ImageIDPath,
) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.DeleteImage")
	defer span.End()

	if err := h.imageSvc.Delete(ctx, imageID); err != nil {
		gen.RespondError(w, r, fmt.Errorf("deleting image: %w", err))
		return
	}

	gen.RespondNoContent(w, http.StatusOK)
}
