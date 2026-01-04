package images

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/pkg/apperr"
)

type Fit string

const (
	// Cover crops the image to completely fill the given dimensions.
	FitCover Fit = "COVER"

	// Contain resizes the image to fit within the given dimensions while
	// maintaining the aspect ratio. Empty areas are filled with background
	// color.
	FitContain Fit = "CONTAIN"

	// Fill resizes the image to the given dimensions without maintaining
	// the aspect ratio. This may distort the image.
	FitFill Fit = "FILL"
)

// Ensure interfaces are implemented
var (
	_ driver.Valuer = Fit("")
	_ sql.Scanner   = (*Fit)(nil)
)

func (f *Fit) GetOrDefault() Fit {
	return lo.FromPtrOr(f, FitCover)
}

func (f Fit) Validate() error {
	switch f {
	case FitCover:
	case FitContain:
	case FitFill:
	default:
		return apperr.NewError(apperr.CodeBadRequest).WithSummary("Unexpected image crop mode %q", f)
	}
	return nil
}

func (f Fit) Value() (driver.Value, error) {
	return string(f), nil
}

func (f *Fit) Scan(value any) error {
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
			WithSummary("Invalid value of image crop mode: %[1]T(%[1]v)", value)
	}

	*f = Fit(str)
	return nil
}
