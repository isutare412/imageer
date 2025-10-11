package port

import (
	"context"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type AuthService interface {
	StartGoogleSignIn(context.Context, domain.StartGoogleSignInRequest) (domain.StartGoogleSignInResponse, error)
	FinishGoogleSignIn(context.Context, domain.FinishGoogleSignInRequest) (domain.FinishGoogleSignInResponse, error)
}
