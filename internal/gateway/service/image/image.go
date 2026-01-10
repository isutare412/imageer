package image

import (
	"fmt"
	"regexp"

	"github.com/isutare412/imageer/pkg/images"
)

var imageS3KeyPattern = regexp.MustCompile(`projects/([^/]+)/images/([^/]+)`)

func imageBasePath(projectID, imageID string) string {
	return fmt.Sprintf("projects/%s/images/%s", projectID, imageID)
}

func (s *Service) imageS3BasePath(projectID, imageID string) string {
	return fmt.Sprintf("%s/%s", s.cfg.S3KeyPrefix, imageBasePath(projectID, imageID))
}

func (s *Service) imageS3Key(projectID, imageID string, format images.Format) string {
	base := s.imageS3BasePath(projectID, imageID)
	return fmt.Sprintf("%s/original.%s", base, format.Extension())
}

func (s *Service) imagePublicURL(projectID, imageID string, format images.Format) string {
	return fmt.Sprintf("%s/%s/original.%s", s.cfg.CDNDomain, imageBasePath(projectID, imageID),
		format.Extension())
}

func (s *Service) imageVariantS3Key(projectID, imageID, presetID string, format images.Format,
) string {
	base := s.imageS3BasePath(projectID, imageID)
	return fmt.Sprintf("%s/variants/%s.%s", base, presetID, format.Extension())
}

func (s *Service) imageVariantPublicURL(projectID, imageID, presetID string, format images.Format,
) string {
	return fmt.Sprintf("%s/%s/variants/%s.%s", s.cfg.CDNDomain, imageBasePath(projectID, imageID),
		presetID, format.Extension())
}

func parseImageS3Key(key string) (projectID, imageID string, ok bool) {
	matches := imageS3KeyPattern.FindStringSubmatch(key)
	if len(matches) != 3 {
		return "", "", false
	}
	return matches[1], matches[2], true
}
