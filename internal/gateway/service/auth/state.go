package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/apperr"
)

func (s *AuthService) createOIDCState(r *http.Request) (string, error) {
	state := domain.OIDCState{
		OriginURL: httpReferer(r),
	}

	stateBytes, err := json.Marshal(state)
	if err != nil {
		return "", apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("failed to marshal OIDC state").
			WithCause(err)
	}

	stateEcrypted, err := s.crypter.Encrypt(stateBytes)

	return base64.RawURLEncoding.EncodeToString(stateEcrypted), nil
}

func (s *AuthService) decryptOIDCState(state string) (domain.OIDCState, error) {
	stateBytes, err := base64.RawURLEncoding.DecodeString(state)
	if err != nil {
		return domain.OIDCState{}, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("unexpected OIDC state format").
			WithCause(err)
	}

	stateDecrypted, err := s.crypter.Decrypt(stateBytes)
	if err != nil {
		return domain.OIDCState{}, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("cannot verify OIDC state").
			WithCause(err)
	}

	var oidcState domain.OIDCState
	if err := json.Unmarshal(stateDecrypted, &oidcState); err != nil {
		return domain.OIDCState{}, apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("failed to unmarshal OIDC state").
			WithCause(err)
	}

	return oidcState, nil
}
