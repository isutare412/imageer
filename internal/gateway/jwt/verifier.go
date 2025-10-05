package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/apperr"
)

type Verifier struct {
	keyChain rsaKeyChain
}

func NewVerifier(cfg VerifierConfig) (*Verifier, error) {
	keyChain, err := newRSAKeyChain(cfg.KeyPairs)
	if err != nil {
		return nil, fmt.Errorf("creating RSA key chain: %w", err)
	}

	return &Verifier{
		keyChain: keyChain,
	}, nil
}

func (v *Verifier) VerifyUserToken(token string) (payload domain.UserTokenPayload, err error) {
	jwtToken, err := jwt.ParseWithClaims(token, &appClaims{}, v.publicKeyPicker(),
		jwt.WithIssuedAt(), jwt.WithIssuer(imageerGatewayIssuer))
	switch {
	case err != nil:
		return payload, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("failed to parse and verify JWT token").
			WithCause(err)
	case !jwtToken.Valid:
		return payload, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("invalid token")
	}

	claims, ok := jwtToken.Claims.(*appClaims)
	if !ok {
		return payload, apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("unexpected claims type")
	}

	return claims.toUserTokenPayload(), nil
}

func (v *Verifier) publicKeyPicker() func(jwtToken *jwt.Token) (any, error) {
	return func(jwtToken *jwt.Token) (any, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, apperr.NewError(apperr.CodeBadRequest).
				WithSummary("unexpected signing method: %v", jwtToken.Header["alg"])
		}

		kid, ok := jwtToken.Header["kid"].(string)
		if !ok {
			return nil, apperr.NewError(apperr.CodeBadRequest).
				WithSummary("kid of token is not found or not string")
		}

		pair, ok := v.keyChain[kid]
		if !ok {
			return nil, apperr.NewError(apperr.CodeBadRequest).
				WithSummary("kid %s is not found from key chain", kid)
		}

		return pair.public, nil
	}
}
