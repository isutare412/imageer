package port

import "context"

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type HealthChecker interface {
	ComponentName() string
	HealthCheck(context.Context) error
}
