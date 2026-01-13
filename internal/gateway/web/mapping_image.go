package web

import (
	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

func ImageToWeb(img domain.Image) Image {
	return Image{
		ID:        img.ID,
		CreatedAt: img.CreatedAt,
		UpdatedAt: img.UpdatedAt,
		Format:    img.Format,
		State:     img.State,
		URL:       img.URL,
		Variants: lo.Map(img.Variants, func(iv domain.ImageVariant, _ int) ImageVariant {
			return ImageVariantToWeb(iv)
		}),
	}
}

func ImageVariantToWeb(iv domain.ImageVariant) ImageVariant {
	return ImageVariant{
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

func CreateUploadURLRequestToDomain(projID string, req CreateUploadURLRequest,
) domain.CreateUploadURLRequest {
	return domain.CreateUploadURLRequest{
		ProjectID:   projID,
		FileName:    req.FileName,
		Format:      req.Format,
		PresetNames: req.PresetNames,
	}
}

func UploadURLToWeb(u domain.UploadURL) UploadURL {
	return UploadURL{
		ImageID:   u.ImageID,
		ExpiresAt: u.ExpiresAt,
		URL:       u.URL,
		Header: lo.MapEntries(u.Header, func(k string, v []string) (string, string) {
			return k, v[0]
		}),
	}
}
