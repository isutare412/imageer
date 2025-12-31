package users

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/isutare412/imageer/pkg/apperr"
)

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleGuest Role = "GUEST"
)

// Ensure interfaces are implemented
var (
	_ driver.Valuer = Role("")
	_ sql.Scanner   = (*Role)(nil)
)

func (s Role) Validate() error {
	switch s {
	case RoleAdmin:
	case RoleGuest:
	default:
		return apperr.NewError(apperr.CodeBadRequest).WithSummary("Unexpected user role %q", s)
	}
	return nil
}

func (s Role) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *Role) Scan(value any) error {
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
			WithSummary("Invalid value of user role: %[1]T(%[1]v)", value)
	}

	*s = Role(str)
	return nil
}
