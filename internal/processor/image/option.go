package image

import (
	"github.com/h2non/bimg"

	"github.com/isutare412/imageer/internal/processor/domain"
	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/images"
)

func imageTypeToFormat(t string) (images.Format, error) {
	switch t {
	case "jpeg", "jpg":
		return images.FormatJPEG, nil
	case "png":
		return images.FormatPNG, nil
	case "webp":
		return images.FormatWebp, nil
	case "avif":
		return images.FormatAVIF, nil
	case "heif", "heic":
		return images.FormatHEIC, nil
	default:
		return images.Format(""), apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Unexpected image type %s", t)
	}
}

func applyPreset(o *bimg.Options, t domain.Preset) {
	o.StripMetadata = true
	o.Quality = int(t.Quality)

	var typ bimg.ImageType
	switch t.Format {
	case images.FormatJPEG:
		typ = bimg.JPEG
	case images.FormatPNG:
		typ = bimg.PNG
	case images.FormatWebp:
		typ = bimg.WEBP
	case images.FormatAVIF:
		typ = bimg.AVIF
	default:
		typ = bimg.WEBP
	}
	o.Type = typ

	if t.Width != nil {
		o.Width = int(*t.Width)
	}
	if t.Height != nil {
		o.Height = int(*t.Height)
	}

	if t.Fit != nil {
		switch *t.Fit {
		case images.FitCover:
			o.Crop = true
		case images.FitContain:
			o.Embed = true
		case images.FitFill:
			o.Force = true
		}
	}

	if t.Anchor != nil {
		var anchor bimg.Gravity
		switch *t.Anchor {
		case images.AnchorSmart:
			anchor = bimg.GravitySmart
		case images.AnchorCenter:
			anchor = bimg.GravityCentre
		case images.AnchorNorth:
			anchor = bimg.GravityNorth
		case images.AnchorSouth:
			anchor = bimg.GravitySouth
		case images.AnchorWest:
			anchor = bimg.GravityWest
		case images.AnchorEast:
			anchor = bimg.GravityEast
		}
		o.Gravity = anchor
	}
}
