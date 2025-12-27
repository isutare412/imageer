package log

//go:generate go tool enumer -type=Format -trimprefix Format -output format_enum.go -transform lower -text
type Format int

const (
	FormatJSON Format = iota
	FormatText
	FormatPretty
)
