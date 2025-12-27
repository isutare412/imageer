package dbhelpers

//go:generate go tool enumer -type=SortDirection  -trimprefix SortDirection -output sort_direction_enum.go -transform snake-upper -text -json -sql
type SortDirection int

const (
	SortDirectionUnspecified SortDirection = iota
	SortDirectionAsc
	SortDirectionDesc
)
