package serviceaccounts

//go:generate go tool enumer -type=AccessScope -trimprefix AccessScope -output access_scope_enum.go -transform snake-upper -text -json -sql
type AccessScope int

const (
	AccessScopeFull AccessScope = iota
	AccessScopeProject
)
