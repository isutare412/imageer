package web

import "github.com/isutare412/imageer/internal/gateway/domain"

func TransformationToWeb(t domain.Transformation) Transformation {
	return Transformation{
		ID:        t.ID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		Name:      t.Name,
		Default:   t.Default,
		Width:     t.Width,
		Height:    t.Height,
	}
}

func CreateTransformationRequestToDomain(req CreateTransformationRequest) domain.CreateTransformationRequest {
	return domain.CreateTransformationRequest{
		Name:    req.Name,
		Default: req.Default,
		Width:   req.Width,
		Height:  req.Height,
	}
}

func UpsertTransformationRequestToDomain(req UpsertTransformationRequest) domain.UpsertTransformationRequest {
	return domain.UpsertTransformationRequest{
		ID:      req.ID,
		Name:    req.Name,
		Default: req.Default,
		Width:   req.Width,
		Height:  req.Height,
	}
}
