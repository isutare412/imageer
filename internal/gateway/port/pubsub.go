package port

import "context"

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type ImageProcessDonePublisher interface {
	Publish(ctx context.Context, imageID string) (receiveCount int64, err error)
}

type ImageProcessDoneSubscriber interface {
	// Subscribe returns a channel that emits struct{}{} on each notification.
	// Channel closes when context is cancelled or on error.
	Subscribe(ctx context.Context, imageID string) (<-chan struct{}, <-chan error)
}
