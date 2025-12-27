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
	ContentType      images.ContentType
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
	TransformationID   string
	TransformationName string
	URL                string
}

type PresignedURL struct {
	ImageID   string
	ExpiresAt time.Time
	UploadURL string
}

type CreatePresignedURLRequest struct {
	FileName            string             `validate:"required,max=1024"`
	ContentType         images.ContentType `validate:"validateFn=IsAContentType"`
	TransformationNames []string           `validate:"dive,required,max=64"`
}
