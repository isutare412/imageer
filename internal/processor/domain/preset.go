package domain

import (
	"github.com/samber/lo"

	"github.com/isutare412/imageer/pkg/images"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
)

type Preset struct {
	ID      string
	Name    string
	Default bool

	Format  images.Format
	Quality images.Quality
	Fit     *images.Fit
	Anchor  *images.Anchor
	Width   *int32
	Height  *int32
}

func NewPreset(p *imageerv1.Preset) Preset {
	return Preset{
		ID:      p.Id,
		Name:    p.Name,
		Default: p.Default,
		Format:  images.NewFormatFromProto(p.Format),
		Quality: images.Quality(p.Quality),
		Fit:     lo.EmptyableToPtr(images.NewFitFromProto(p.Fit)),
		Anchor:  lo.EmptyableToPtr(images.NewAnchorFromProto(p.Anchor)),
		Width:   p.Width,
		Height:  p.Height,
	}
}
