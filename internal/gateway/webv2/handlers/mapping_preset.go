package handlers

import (
	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
	"github.com/isutare412/imageer/pkg/images"
)

func PresetToWeb(t domain.Preset) gen.Preset {
	return gen.Preset{
		ID:        t.ID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		Name:      t.Name,
		Default:   t.Default,
		Format:    t.Format,
		Quality:   int64(t.Quality),
		Fit:       t.Fit,
		Anchor:    t.Anchor,
		Width:     t.Width,
		Height:    t.Height,
	}
}

func CreatePresetRequestToDomain(req gen.CreatePresetRequest) domain.CreatePresetRequest {
	var quality *images.Quality
	if req.Quality != nil {
		q := images.Quality(*req.Quality)
		quality = &q
	}

	return domain.CreatePresetRequest{
		Name:    req.Name,
		Default: req.Default,
		Format:  req.Format,
		Quality: quality,
		Fit:     req.Fit,
		Anchor:  req.Anchor,
		Width:   req.Width,
		Height:  req.Height,
	}
}

func UpsertPresetRequestToDomain(req gen.UpsertPresetRequest) domain.UpsertPresetRequest {
	var quality *images.Quality
	if req.Quality != nil {
		q := images.Quality(*req.Quality)
		quality = &q
	}

	return domain.UpsertPresetRequest{
		ID:      req.ID,
		Name:    req.Name,
		Default: req.Default,
		Format:  req.Format,
		Quality: quality,
		Fit:     req.Fit,
		Anchor:  req.Anchor,
		Width:   req.Width,
		Height:  req.Height,
	}
}
