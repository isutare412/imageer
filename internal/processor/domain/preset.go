package domain

import (
	"github.com/isutare412/imageer/pkg/images"
)

type Preset struct {
	ID      string
	Name    string
	Default bool

	Format  images.Format
	Quality images.Quality
	Fit     *images.Fit
	Anchor  *images.Anchor
	Width   *int64
	Height  *int64
}
