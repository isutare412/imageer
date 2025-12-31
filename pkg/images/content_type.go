package images

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/isutare412/imageer/pkg/apperr"
)

type ContentType string

const (
	ContentTypeImageJPEG ContentType = "image/jpeg"
	ContentTypeImagePNG  ContentType = "image/png"
	ContentTypeImageWebp ContentType = "image/webp"
)

// Ensure interfaces are implemented
var (
	_ driver.Valuer = ContentType("")
	_ sql.Scanner   = (*ContentType)(nil)
)

func (t ContentType) Validate() error {
	switch t {
	case ContentTypeImageJPEG:
	case ContentTypeImagePNG:
	case ContentTypeImageWebp:
	default:
		return apperr.NewError(apperr.CodeBadRequest).WithSummary("Unexpected image content type %q", t)
	}
	return nil
}

func (t ContentType) Value() (driver.Value, error) {
	return string(t), nil
}

func (t *ContentType) Scan(value any) error {
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
			WithSummary("Invalid value of image content type: %[1]T(%[1]v)", value)
	}

	*t = ContentType(str)
	return nil
}
