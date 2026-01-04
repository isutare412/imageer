package images

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/pkg/apperr"
)

type Anchor string

const (
	AnchorSmart  Anchor = "SMART"
	AnchorCenter Anchor = "CENTER"
	AnchorNorth  Anchor = "NORTH"
	AnchorSouth  Anchor = "SOUTH"
	AnchorEast   Anchor = "EAST"
	AnchorWest   Anchor = "WEST"
)

// Ensure interfaces are implemented
var (
	_ driver.Valuer = Anchor("")
	_ sql.Scanner   = (*Anchor)(nil)
)

func (a *Anchor) GetOrDefault() Anchor {
	return lo.FromPtrOr(a, AnchorSmart)
}

func (a Anchor) Validate() error {
	switch a {
	case AnchorSmart:
	case AnchorCenter:
	case AnchorNorth:
	case AnchorSouth:
	case AnchorEast:
	case AnchorWest:
	default:
		return apperr.NewError(apperr.CodeBadRequest).WithSummary("Unexpected image anchor %q", a)
	}
	return nil
}

func (a Anchor) Value() (driver.Value, error) {
	return string(a), nil
}

func (a *Anchor) Scan(value any) error {
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
			WithSummary("Invalid value of image anchor: %[1]T(%[1]v)", value)
	}

	*a = Anchor(str)
	return nil
}
