package images

type State string

const (
	StateWaitingUpload State = "WAITING_UPLOAD"
	StateProcessing    State = "PROCESSING"
	StateReady         State = "READY"
)

type ContentType string

const (
	ContentTypeJPEG ContentType = "image/jpeg"
	ContentTypePNG  ContentType = "image/png"
	ContentTypeWEBP ContentType = "image/webp"
)
