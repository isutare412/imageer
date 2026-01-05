package images

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/pkg/apperr"
)

type Format string

const (
	FormatJPEG Format = "JPEG"
	FormatPNG  Format = "PNG"
	FormatWebp Format = "WEBP"
	FormatAVIF Format = "AVIF"
	FormatHEIC Format = "HEIC"
)

// Ensure interfaces are implemented
var (
	_ driver.Valuer = Format("")
	_ sql.Scanner   = (*Format)(nil)
)

func (f *Format) GetOrDefault() Format {
	return lo.FromPtrOr(f, FormatWebp)
}

func (f Format) Validate() error {
	switch f {
	case FormatJPEG:
	case FormatPNG:
	case FormatWebp:
	case FormatAVIF:
	case FormatHEIC:
	default:
		return apperr.NewError(apperr.CodeBadRequest).WithSummary("Unexpected image format %q", f)
	}
	return nil
}

func (f Format) ValidateForPreset() error {
	switch f {
	case FormatJPEG:
	case FormatPNG:
	case FormatWebp:
	default:
		return apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Unexpected image format %q for preset", f)
	}
	return nil
}

func (f Format) ToExtension() string {
	switch f {
	case FormatJPEG:
		return "jpg"
	case FormatPNG:
		return "png"
	case FormatWebp:
		return "webp"
	case FormatAVIF:
		return "avif"
	case FormatHEIC:
		return "heif"
	default:
		return ""
	}
}

func (f Format) Value() (driver.Value, error) {
	return string(f), nil
}

func (f *Format) Scan(value any) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("Invalid value of image format: %[1]T(%[1]v)", value)
	}

	*f = Format(str)
	return nil
}
