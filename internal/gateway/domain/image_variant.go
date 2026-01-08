package domain

import (
	"time"

	"github.com/isutare412/imageer/pkg/images"
)

type ImageVariant struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Format    images.Format
	State     images.VariantState
	S3Key     string
	URL       string

	ImageID string
	Preset  PresetReference
}
