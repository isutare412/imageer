package jwt

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"

	"github.com/isutare412/imageer/pkg/apperr"
)

type rsaKeyPair struct {
	name    string
	private *rsa.PrivateKey
	public  *rsa.PublicKey
}

// rsaKeyChain maps key pair name to rsaKeyPair.
type rsaKeyChain map[string]rsaKeyPair

func newRSAKeyChain(pairs []RSAKeyBytesPair) (rsaKeyChain, error) {
	chain := make(rsaKeyChain, len(pairs))
	for _, pair := range pairs {
		prv, err := jwt.ParseRSAPrivateKeyFromPEM(pair.Private)
		if err != nil {
			return nil, apperr.NewError(apperr.CodeInternalServerError).
				WithSummary("failed to parse RSA private key").
				WithCause(err)
		}

		pub, err := jwt.ParseRSAPublicKeyFromPEM(pair.Public)
		if err != nil {
			return nil, apperr.NewError(apperr.CodeInternalServerError).
				WithSummary("failed to parse RSA public key").
				WithCause(err)
		}

		chain[pair.Name] = rsaKeyPair{
			name:    pair.Name,
			private: prv,
			public:  pub,
		}
	}

	return chain, nil
}
