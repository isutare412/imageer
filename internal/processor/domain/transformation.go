package domain

import (
	"github.com/isutare412/imageer/pkg/images"
)

type Transformation struct {
	ID      string
	Name    string
	Default bool

	Format  images.Format
	Quality images.Quality
	Fit     *images.Fit
	Width   *int64
	Height  *int64

	Crop   bool
	Anchor *images.Anchor
}
