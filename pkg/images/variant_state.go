package images

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/isutare412/imageer/pkg/apperr"
)

type VariantState string

const (
	VariantStateWaitingUpload VariantState = "WAITING_UPLOAD"
	VariantStateProcessing    VariantState = "PROCESSING"
	VariantStateFailed        VariantState = "FAILED"
	VariantStateReady         VariantState = "READY"
)

// Ensure interfaces are implemented
var (
	_ driver.Valuer = VariantState("")
	_ sql.Scanner   = (*VariantState)(nil)
)

func (s VariantState) Validate() error {
	switch s {
	case VariantStateProcessing:
	case VariantStateFailed:
	case VariantStateReady:
	default:
		return apperr.NewError(apperr.CodeBadRequest).WithSummary("Unexpected image variant state %q", s)
	}
	return nil
}

func (s VariantState) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *VariantState) Scan(value any) error {
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
			WithSummary("Invalid value of image variant state: %[1]T(%[1]v)", value)
	}

	*s = VariantState(str)
	return nil
}
