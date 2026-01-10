package domain

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/isutare412/imageer/pkg/images"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
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

func (i ImageVariant) ToProto() *imageerv1.ImageVariant {
	return &imageerv1.ImageVariant{
		Id:        i.ID,
		CreatedAt: timestamppb.New(i.CreatedAt),
		UpdatedAt: timestamppb.New(i.UpdatedAt),
		Format:    i.Format.ToProto(),
		State:     i.State.ToProto(),
		S3Key:     i.S3Key,
		Url:       i.URL,
		ImageId:   i.ImageID,
	}
}
