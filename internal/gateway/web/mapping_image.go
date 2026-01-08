package web

import (
	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

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
