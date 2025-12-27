package images

//go:generate go tool enumer -type=State -trimprefix State -output state_enum.go -transform snake-upper -text -json -sql
type State int

const (
	StateWaitingUpload State = iota
	StateProcessing
	StateReady
)
