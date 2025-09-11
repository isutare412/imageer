package oidc

import (
	"github.com/coreos/go-oidc/v3/oidc"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/apperr"
)

// googleIDTokenClaim represents the structure of a Google ID token.
// ref: https://developers.google.com/identity/openid-connect/openid-connect#obtainuserinfo
type googleIDTokenClaim struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Name          string `json:"name"`
	PictureURL    string `json:"picture"`
	ProfileURL    string `json:"profile"`
}

func googleIDTokenToDomain(token *oidc.IDToken) (payload domain.IDTokenPayload, err error) {
	var claims googleIDTokenClaim
	if err := token.Claims(&claims); err != nil {
		return payload, apperr.NewError(apperr.CodeBadRequest).
			WithSummary("extracting claims from ID token").
			WithCause(err)
	}

	return domain.IDTokenPayload{
		Audience:      token.Audience[0],
		Issuer:        token.Issuer,
		Subject:       token.Subject,
		Email:         claims.Email,
		EmailVerified: &claims.EmailVerified,
		FamilyName:    &claims.FamilyName,
		GivenName:     &claims.GivenName,
		FullName:      claims.Name,
		PictureURL:    &claims.PictureURL,
		ProfileURL:    &claims.ProfileURL,
		Expiry:        token.Expiry,
		IssuedAt:      token.IssuedAt,
	}, nil
}
