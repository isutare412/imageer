package port

import "github.com/isutare412/imageer/internal/gateway/domain"

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type JWTSigner interface {
	SignUserToken(payload domain.UserTokenPayload) (token string, err error)
}

type JWTVerifier interface {
	VerifyUserToken(token string) (payload domain.UserTokenPayload, err error)
}
