package port

import (
	"context"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type UserRepository interface {
	FindByID(ctx context.Context, id string) (user domain.User, err error)
	Upsert(ctx context.Context, user domain.User) (userCreated domain.User, err error)
}
