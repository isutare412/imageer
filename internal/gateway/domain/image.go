package domain

import (
	"net/http"
	"time"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/isutare412/imageer/pkg/dbhelpers"
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

func (i Image) AllVariantsProcessed() bool {
	if len(i.Variants) == 0 {
		return true
	}
	for _, v := range i.Variants {
		if !v.State.IsTerminal() {
			return false
		}
	}
	return true
}

type Images struct {
	Items []Image
	Total int64
}

type ListImagesParams struct {
	Offset *int `validate:"omitempty,min=0"`
	Limit  *int `validate:"omitempty,min=1,max=100"`

	SearchFilter ImageSearchFilter
	SortFilter   ImageSortFilter
}

func (p ListImagesParams) OffsetOrDefault() int {
	return lo.FromPtrOr(p.Offset, 0)
}

func (p ListImagesParams) LimitOrDefault() int {
	return lo.FromPtrOr(p.Limit, 20)
}

type ImageSearchFilter struct {
	ProjectID       *string
	State           *images.State
	UpdatedAtBefore *time.Time
}

type ImageSortFilter struct {
	CreatedAt bool
	UpdatedAt bool
	Direction dbhelpers.SortDirection
}

type UpdateImageRequest struct {
	ID    string
	State *images.State
}

type ReprocessImagesRequest struct {
	ImageIDs     []string
	ReprocessAll bool
}

type ImageProcessingLog struct {
	ID             int
	CreatedAt      time.Time
	IsSuccess      bool
	ErrorCode      *int
	ErrorMessage   *string
	ElapsedTime    time.Duration
	ImageVariantID string
}

func NewImageProcessingLog(res *imageerv1.ImageProcessResult) ImageProcessingLog {
	return ImageProcessingLog{
		IsSuccess:      res.IsSuccess,
		ErrorCode:      lo.EmptyableToPtr(int(res.ErrorCode)),
		ErrorMessage:   lo.EmptyableToPtr(res.ErrorMessage),
		ElapsedTime:    res.ProcessingTime.AsDuration(),
		ImageVariantID: res.ImageVariantId,
	}
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
