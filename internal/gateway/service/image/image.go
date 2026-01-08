package image

import (
	"fmt"

	"github.com/isutare412/imageer/pkg/images"
)

func (s *Service) imageS3BasePath(projectID, imageID string) string {
	return fmt.Sprintf("%s/projects/%s/images/%s", s.cfg.S3KeyPrefix, projectID, imageID)
}

func (s *Service) imageS3Key(projectID, imageID string, format images.Format) string {
	base := s.imageS3BasePath(projectID, imageID)
	return fmt.Sprintf("%s.%s", base, format.Extension())
}

func (s *Service) imageVariantS3Key(projectID, imageID, presetID string, format images.Format) string {
	base := s.imageS3BasePath(projectID, imageID)
	return fmt.Sprintf("%s/variants/%s.%s", base, presetID, format.Extension())
}
