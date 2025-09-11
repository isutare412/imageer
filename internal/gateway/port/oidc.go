package port

import (
	"context"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type OIDCProvider interface {
	BuildAuthenticationURL(baseURL, state string) string
	ExchangeCode(ctx context.Context, baseURL, code string) (domain.IDTokenPayload, error)
}
