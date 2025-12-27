package images

//go:generate go tool enumer -type=ContentType -trimprefix ContentType -output content_type_enum.go -transform snake-upper -text -json -sql
type ContentType int

const (
	ContentTypeImageJPEG ContentType = iota
	ContentTypeImagePNG
	ContentTypeImageWEBP
)
