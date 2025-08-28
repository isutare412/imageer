package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/isutare412/imageer/pkg/images"
)

type Image struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	State     images.State
	URLSet    ImageURLSet
}

type Images struct {
	Items []Image
	Total int64
}

type ReprocessImagesRequest struct {
	ImageIDs     []uuid.UUID
	ReprocessAll bool
}

type Transformation struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Default   bool
	Width     int64
	Height    int64
}

type CreateTransformationRequest struct {
	Name    string `validate:"required,max=64"`
	Default bool
	Width   int64 `validate:"min=1"`
	Height  int64 `validate:"min=1"`
}

type UpdateTransformationRequest struct {
	ID      uuid.UUID `validate:"required"`
	Name    *string   `validate:"omitempty,max=64"`
	Default *bool
	Width   *int64 `validate:"omitempty,min=1"`
	Height  *int64 `validate:"omitempty,min=1"`
}

type DeleteTransformationRequest struct {
	ID uuid.UUID `validate:"required"`
}

type ImageURLSet struct {
	OriginalURL     string
	Transformations []VariantURL
}

type VariantURL struct {
	TransformationID   uuid.UUID
	TransformationName string
	URL                string
}

type PresignedURL struct {
	ImageID   uuid.UUID
	ExpiresAt time.Time
	UploadURL string
}

type CreatePresignedURLRequest struct {
	FileName            string             `validate:"required,max=1024"`
	ContentType         images.ContentType `validate:"oneof=image/jpeg image/png image/webp"`
	TransformationNames []string           `validate:"dive,required,max=64"`
}
