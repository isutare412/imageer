package serviceaccounts

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/isutare412/imageer/pkg/apperr"
)

type AccessScope string

const (
	AccessScopeFull    AccessScope = "FULL"
	AccessScopeProject AccessScope = "PROJECT"
)

// Ensure interfaces are implemented
var (
	_ driver.Valuer = AccessScope("")
	_ sql.Scanner   = (*AccessScope)(nil)
)

func (s AccessScope) Validate() error {
	switch s {
	case AccessScopeFull:
	case AccessScopeProject:
	default:
		return apperr.NewError(apperr.CodeBadRequest).
			WithSummary("Unexpected service account access scope %q", s)
	}
	return nil
}

func (s AccessScope) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *AccessScope) Scan(value any) error {
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
			WithSummary("Invalid value of access scope access scope: %[1]T(%[1]v)", value)
	}

	*s = AccessScope(str)
	return nil
}
