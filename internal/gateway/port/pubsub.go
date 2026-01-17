package port

import "context"

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

// ImageNotificationPublisher publishes image-related events (upload done, process done).
type ImageNotificationPublisher interface {
	PublishUploadDone(ctx context.Context, imageID string) (receiveCount int64, err error)
	PublishProcessDone(ctx context.Context, imageID string) (receiveCount int64, err error)
}

// ImageUploadDoneSubscriber subscribes to image upload done notifications.
type ImageUploadDoneSubscriber interface {
	// Subscribe returns a channel that emits struct{}{} on each notification.
	// Channel closes when context is cancelled or on error.
	Subscribe(ctx context.Context, imageID string) (<-chan struct{}, <-chan error)
}

// ImageProcessDoneSubscriber subscribes to image process done notifications.
type ImageProcessDoneSubscriber interface {
	// Subscribe returns a channel that emits struct{}{} on each notification.
	// Channel closes when context is cancelled or on error.
	Subscribe(ctx context.Context, imageID string) (<-chan struct{}, <-chan error)
}
