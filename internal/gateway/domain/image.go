package domain

import (
	"net/http"
	"time"

	"github.com/isutare412/imageer/pkg/images"
)

type Image struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	FileName  string
	Format    images.Format
	State     images.State
	S3Key     string
	URL       string
	Variants  []ImageVariant
	Project   ProjectReference
}

type Images struct {
	Items []Image
	Total int64
}

type ReprocessImagesRequest struct {
	ImageIDs     []string
	ReprocessAll bool
}

type UploadURL struct {
	ImageID   string
	ExpiresAt time.Time
	URL       string
	Header    http.Header
}

type CreateUploadURLRequest struct {
	FileName    string        `validate:"required,max=512"`
	Format      images.Format `validate:"validateFn=Validate"`
	PresetNames []string      `validate:"dive,required,max=64"`
}

type PresignPutObjectRequest struct {
	S3Key       string
	ContentType string
}

type PresignPutObjectResponse struct {
	URL      string
	Header   http.Header
	ExpireAt time.Time
}
