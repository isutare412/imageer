package port

import "context"

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type ImageProcessDonePublisher interface {
	Publish(ctx context.Context, imageID string) (receiveCount int64, err error)
}

type ImageProcessDoneSubscriber interface {
	Wait(ctx context.Context, imageID string) error
}
