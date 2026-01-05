package domain

import (
	"time"

	"github.com/isutare412/imageer/pkg/images"
)

type Image struct {
	ID               string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	OriginalFileName string
	Format           images.Format
	State            images.State
	URLSet           ImageURLSet
}

type Images struct {
	Items []Image
	Total int64
}

type ReprocessImagesRequest struct {
	ImageIDs     []string
	ReprocessAll bool
}

type ImageURLSet struct {
	OriginalURL string
	Variants    []VariantURL
}

type VariantURL struct {
	PresetID   string
	PresetName string
	URL        string
}

type UploadURL struct {
	ImageID   string
	ExpiresAt time.Time
	URL       string
}

type CreateUploadURLRequest struct {
	FileName    string        `validate:"required,max=1024"`
	Format      images.Format `validate:"validateFn=Validate"`
	PresetNames []string      `validate:"dive,required,max=64"`
}
