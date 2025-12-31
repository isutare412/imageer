package log

import "github.com/isutare412/imageer/pkg/apperr"

type Format string

const (
	FormatJSON   Format = "json"
	FormatText   Format = "text"
	FormatPretty Format = "pretty"
)

func (f Format) Validate() error {
	switch f {
	case FormatJSON:
	case FormatText:
	case FormatPretty:
	default:
		return apperr.NewError(apperr.CodeBadRequest).WithSummary("Unexpected log format %q", f)
	}
	return nil
}
