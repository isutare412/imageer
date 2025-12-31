package images

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/isutare412/imageer/pkg/apperr"
)

type State string

const (
	StateWaitingUpload State = "WAITING_UPLOAD"
	StateProcessing    State = "PROCESSING"
	StateReady         State = "READY"
)

// Ensure interfaces are implemented
var (
	_ driver.Valuer = State("")
	_ sql.Scanner   = (*State)(nil)
)

func (s State) Validate() error {
	switch s {
	case StateWaitingUpload:
	case StateProcessing:
	case StateReady:
	default:
		return apperr.NewError(apperr.CodeBadRequest).WithSummary("Unexpected image state %q", s)
	}
	return nil
}

func (s State) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *State) Scan(value any) error {
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
			WithSummary("Invalid value of image state: %[1]T(%[1]v)", value)
	}

	*s = State(str)
	return nil
}
