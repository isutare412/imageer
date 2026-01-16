package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/apperr"
)

func (s *Service) createOIDCState(r *http.Request, redirectPath string) (string, error) {
	state := domain.OIDCState{
		RedirectURL: httpRedirectURL(r, redirectPath),
	}

	stateBytes, err := json.Marshal(state)
	if err != nil {
		return "", apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("failed to marshal OIDC state").
			WithCause(err)
	}

	stateEcrypted, err := s.crypter.Encrypt(stateBytes)
	if err != nil {
		return "", fmt.Errorf("ecrypting OIDC state: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(stateEcrypted), nil
}

func (s *Service) decryptOIDCState(state string) (oidcState domain.OIDCState, err error) {
	stateBytes, err := base64.RawURLEncoding.DecodeString(state)
	if err != nil {
		return oidcState, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("unexpected OIDC state format").
			WithCause(err)
	}

	stateDecrypted, err := s.crypter.Decrypt(stateBytes)
	if err != nil {
		return oidcState, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("cannot verify OIDC state").
			WithCause(err)
	}

	if err := json.Unmarshal(stateDecrypted, &oidcState); err != nil {
		return oidcState, apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("failed to unmarshal OIDC state").
			WithCause(err)
	}

	return oidcState, nil
}
