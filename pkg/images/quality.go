package images

import (
	"database/sql"
	"database/sql/driver"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/pkg/apperr"
)

type Quality int

// Ensure interfaces are implemented
var (
	_ driver.Valuer = Quality(0)
	_ sql.Scanner   = (*Quality)(nil)
)

func (q *Quality) GetOrDefault() Quality {
	return lo.FromPtrOr(q, 80)
}

func (q Quality) Validate() error {
	if q < 1 || q > 100 {
		return apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Image quality should be between 1 and 100")
	}
	return nil
}

func (q Quality) Value() (driver.Value, error) {
	return int64(q), nil
}

func (q *Quality) Scan(value any) error {
	if value == nil {
		return nil
	}

	var num int
	switch v := value.(type) {
	case int:
		num = v
	case int32:
		num = int(v)
	case int64:
		num = int(v)
	default:
		return apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("Invalid value of image content type: %[1]T(%[1]v)", value)
	}

	*q = Quality(num)
	return nil
}
