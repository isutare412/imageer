package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/apperr"
)

type Signer struct {
	keyPair rsaKeyPair
}

func NewSigner(cfg SignerConfig) (*Signer, error) {
	keyChain, err := newRSAKeyChain(cfg.KeyPairs)
	if err != nil {
		return nil, fmt.Errorf("creating RSA key chain: %w", err)
	}

	keyPair, ok := keyChain[cfg.ActiveKeyPairName]
	if !ok {
		return nil, apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("active key pair %q not found in the key chain", cfg.ActiveKeyPairName)
	}

	return &Signer{
		keyPair: keyPair,
	}, nil
}

func (s *Signer) SignUserToken(payload domain.UserTokenPayload) (token string, err error) {
	appClaims := newAppClaims(payload)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, appClaims)
	jwtToken.Header["kid"] = s.keyPair.name

	token, err = jwtToken.SignedString(s.keyPair.private)
	if err != nil {
		return "", apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("failed to sign JWT token").
			WithCause(err)
	}

	return token, nil
}
