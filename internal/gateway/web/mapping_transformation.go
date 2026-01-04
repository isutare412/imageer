package web

import (
	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/images"
)

func TransformationToWeb(t domain.Transformation) Transformation {
	return Transformation{
		ID:        t.ID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		Name:      t.Name,
		Default:   t.Default,
		Format:    t.Format,
		Quality:   int64(t.Quality),
		Fit:       t.Fit,
		Width:     t.Width,
		Height:    t.Height,
		Crop:      t.Crop,
		Anchor:    t.Anchor,
	}
}

func CreateTransformationRequestToDomain(req CreateTransformationRequest) domain.CreateTransformationRequest {
	var quality *images.Quality
	if req.Quality != nil {
		q := images.Quality(*req.Quality)
		quality = &q
	}

	return domain.CreateTransformationRequest{
		Name:    req.Name,
		Default: req.Default,
		Format:  req.Format,
		Quality: quality,
		Fit:     req.Fit,
		Width:   req.Width,
		Height:  req.Height,
		Crop:    lo.FromPtr(req.Crop),
		Anchor:  req.Anchor,
	}
}

func UpsertTransformationRequestToDomain(req UpsertTransformationRequest) domain.UpsertTransformationRequest {
	var quality *images.Quality
	if req.Quality != nil {
		q := images.Quality(*req.Quality)
		quality = &q
	}

	return domain.UpsertTransformationRequest{
		ID:      req.ID,
		Name:    req.Name,
		Default: req.Default,
		Format:  req.Format,
		Quality: quality,
		Fit:     req.Fit,
		Width:   req.Width,
		Height:  req.Height,
		Crop:    req.Crop,
		Anchor:  req.Anchor,
	}
}
