package webv2

import (
	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
	"github.com/isutare412/imageer/pkg/dbhelpers"
)

func ImageToWeb(img domain.Image) gen.Image {
	return gen.Image{
		ID:        img.ID,
		CreatedAt: img.CreatedAt,
		UpdatedAt: img.UpdatedAt,
		Format:    img.Format,
		State:     img.State,
		URL:       img.URL,
		Variants: lo.Map(img.Variants, func(iv domain.ImageVariant, _ int) gen.ImageVariant {
			return ImageVariantToWeb(iv)
		}),
	}
}

func ImageVariantToWeb(iv domain.ImageVariant) gen.ImageVariant {
	return gen.ImageVariant{
		ID:         iv.ID,
		CreatedAt:  iv.CreatedAt,
		UpdatedAt:  iv.UpdatedAt,
		Format:     iv.Format,
		State:      iv.State,
		URL:        iv.URL,
		PresetID:   iv.Preset.ID,
		PresetName: iv.Preset.Name,
	}
}

func CreateUploadURLRequestToDomain(projID string, req gen.CreateUploadURLRequest,
) domain.CreateUploadURLRequest {
	return domain.CreateUploadURLRequest{
		ProjectID:   projID,
		FileName:    req.FileName,
		Format:      req.Format,
		PresetNames: req.PresetNames,
	}
}

func UploadURLToWeb(u domain.UploadURL) gen.UploadURL {
	return gen.UploadURL{
		ImageID:   u.ImageID,
		ExpiresAt: u.ExpiresAt,
		URL:       u.URL,
		Header: lo.MapEntries(u.Header, func(k string, v []string) (string, string) {
			return k, v[0]
		}),
	}
}

func ImagesToWeb(imgs domain.Images) gen.Images {
	return gen.Images{
		Items: lo.Map(imgs.Items, func(img domain.Image, _ int) gen.Image {
			return ImageToWeb(img)
		}),
		Total: imgs.Total,
	}
}

func ListImagesAdminParamsToDomain(projectID string, params gen.ListImagesAdminParams,
) domain.ListImagesParams {
	var offset *int
	if params.Offset != nil {
		v := int(*params.Offset)
		offset = &v
	}

	var limit *int
	if params.Limit != nil {
		v := int(*params.Limit)
		limit = &v
	}

	sortFilter := domain.ImageSortFilter{
		Direction: dbhelpers.SortDirectionDesc,
	}
	if params.SortOrder != nil {
		sortFilter.Direction = *params.SortOrder
	}
	if params.SortBy != nil {
		switch *params.SortBy {
		case gen.ListImagesAdminParamsSortByCreatedAt:
			sortFilter.CreatedAt = true
		case gen.ListImagesAdminParamsSortByUpdatedAt:
			sortFilter.UpdatedAt = true
		}
	} else {
		sortFilter.CreatedAt = true
	}

	return domain.ListImagesParams{
		Offset: offset,
		Limit:  limit,
		SearchFilter: domain.ImageSearchFilter{
			ProjectID: &projectID,
		},
		SortFilter: sortFilter,
	}
}

func ListImagesParamsToDomain(projectID string, params gen.ListImagesParams,
) domain.ListImagesParams {
	var offset *int
	if params.Offset != nil {
		v := int(*params.Offset)
		offset = &v
	}

	var limit *int
	if params.Limit != nil {
		v := int(*params.Limit)
		limit = &v
	}

	sortFilter := domain.ImageSortFilter{
		Direction: dbhelpers.SortDirectionDesc,
	}
	if params.SortOrder != nil {
		sortFilter.Direction = *params.SortOrder
	}
	if params.SortBy != nil {
		switch *params.SortBy {
		case gen.CreatedAt:
			sortFilter.CreatedAt = true
		case gen.UpdatedAt:
			sortFilter.UpdatedAt = true
		}
	} else {
		sortFilter.CreatedAt = true
	}

	return domain.ListImagesParams{
		Offset: offset,
		Limit:  limit,
		SearchFilter: domain.ImageSearchFilter{
			ProjectID: &projectID,
		},
		SortFilter: sortFilter,
	}
}
