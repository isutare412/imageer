package domain

import (
	"net/http"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/isutare412/imageer/pkg/images"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
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

func (i Image) ToProto() *imageerv1.Image {
	return &imageerv1.Image{
		Id:        i.ID,
		CreatedAt: timestamppb.New(i.CreatedAt),
		UpdatedAt: timestamppb.New(i.UpdatedAt),
		FileName:  i.FileName,
		Format:    i.Format.ToProto(),
		State:     i.State.ToProto(),
		S3Key:     i.S3Key,
		Url:       i.URL,
		ProjectId: i.Project.ID,
	}
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
	ProjectID   string        `validate:"required,max=36"`
	FileName    string        `validate:"required,max=512"`
	Format      images.Format `validate:"validateFn=Validate"`
	PresetNames []string      `validate:"dive,required,max=64,kebabcase"`
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
