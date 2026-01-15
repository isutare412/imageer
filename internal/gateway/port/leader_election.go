package port

import "context"

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type LeaderHandler interface {
	OnStartedLeading(ctx context.Context)
	OnStoppedLeading()
}
