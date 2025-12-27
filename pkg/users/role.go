package users

//go:generate go tool enumer -type=Role -trimprefix Role -output role_enum.go -transform snake-upper -text -json -sql
type Role int

const (
	RoleAdmin Role = iota
	RoleGuest
)
