package log

import "fmt"

type Format string

const (
	FormatJSON   Format = "json"
	FormatText   Format = "text"
	FormatPretty Format = "pretty"
)

func (f Format) Validate() error {
	switch f {
	case FormatJSON, FormatText, FormatPretty:
		return nil
	default:
		return fmt.Errorf("invalid log format: %q", f)
	}
}
