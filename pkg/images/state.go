package images

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/isutare412/imageer/pkg/apperr"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
)

type State string

const (
	StateUploadPending State = "UPLOAD_PENDING"
	StateUploadExpired State = "UPLOAD_EXPIRED"
	StateFailed        State = "FAILED"
	StateReady         State = "READY"
)

// Ensure interfaces are implemented
var (
	_ driver.Valuer = State("")
	_ sql.Scanner   = (*State)(nil)
)

func (s State) Validate() error {
	switch s {
	case StateUploadPending:
	case StateUploadExpired:
	case StateFailed:
	case StateReady:
	default:
		return apperr.NewError(apperr.CodeBadRequest).WithSummary("Unexpected image state %q", s)
	}
	return nil
}

func (s State) ToProto() imageerv1.ImageState {
	switch s {
	case StateUploadPending:
		return imageerv1.ImageState_IMAGE_STATE_UPLOAD_PENDING
	case StateUploadExpired:
		return imageerv1.ImageState_IMAGE_STATE_UPLOAD_EXPIRED
	case StateFailed:
		return imageerv1.ImageState_IMAGE_STATE_FAILED
	case StateReady:
		return imageerv1.ImageState_IMAGE_STATE_READY
	default:
		return imageerv1.ImageState_IMAGE_STATE_UNSPECIFIED
	}
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
